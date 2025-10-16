"""
FastAPI dependencies
"""
from fastapi import Depends, HTTPException, Header
from typing import Optional
import jwt

from app.config import get_settings
from app.services.feed_service import FeedService
from app.services.ranking_service import RankingService
from app.services.personalization_service import PersonalizationService
from app.ml.recommendation_engine import RecommendationEngine
from app.db import mongodb, redis_client


settings = get_settings()


async def get_current_user(authorization: Optional[str] = Header(None)) -> dict:
    """
    Get current user from JWT token
    """
    if not authorization:
        raise HTTPException(status_code=401, detail="Authorization header missing")
    
    try:
        # Extract token
        scheme, token = authorization.split()
        if scheme.lower() != "bearer":
            raise HTTPException(status_code=401, detail="Invalid authentication scheme")
        
        # Decode JWT
        payload = jwt.decode(
            token,
            settings.JWT_SECRET,
            algorithms=[settings.JWT_ALGORITHM]
        )
        
        user_id = payload.get("user_id")
        if not user_id:
            raise HTTPException(status_code=401, detail="Invalid token payload")
        
        return {"user_id": user_id, "username": payload.get("username", "")}
    
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")
    except Exception as e:
        raise HTTPException(status_code=401, detail=f"Authentication failed: {str(e)}")


def get_recommendation_engine() -> RecommendationEngine:
    """Get recommendation engine instance"""
    return RecommendationEngine()


def get_ranking_service() -> RankingService:
    """Get ranking service instance"""
    return RankingService(
        db=mongodb.get_db(),
        redis_client=redis_client.get_client(),
        elasticsearch_client=None  # Optional
    )


def get_personalization_service() -> PersonalizationService:
    """Get personalization service instance"""
    return PersonalizationService(
        db=mongodb.get_db(),
        redis_client=redis_client.get_client()
    )


def get_feed_service(
    rec_engine: RecommendationEngine = Depends(get_recommendation_engine),
    ranking_service: RankingService = Depends(get_ranking_service),
    personalization_service: PersonalizationService = Depends(get_personalization_service)
) -> FeedService:
    """Get feed service instance"""
    return FeedService(rec_engine, ranking_service, personalization_service)
