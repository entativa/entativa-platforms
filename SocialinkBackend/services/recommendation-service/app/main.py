"""
Socialink Recommendation Service
Main application entry point
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from contextlib import asynccontextmanager
import logging

from app.config import get_settings
from app.api import recommendations


# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Startup and shutdown events
    """
    # Startup
    logger.info("Starting Socialink Recommendation Service...")
    
    # Initialize ML models, DB connections, etc.
    # await initialize_services()
    
    logger.info(f"✅ Recommendation Service ready on port {settings.PORT}")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Recommendation Service...")
    logger.info("👋 Goodbye!")


# Create FastAPI app
app = FastAPI(
    title="Socialink Recommendation Service",
    description="""
    ML-powered recommendation system for Socialink (Facebook-style)
    
    Provides 3 types of recommendations:
    - **Friends Suggestions**: Close friend recommendations (mutual friends priority)
    - **People You May Know**: Broader network recommendations
    - **Communities for You**: Community recommendations
    
    Uses multiple algorithms:
    - Graph-based (friends of friends, social circles)
    - Collaborative filtering (similar users)
    - Content-based (interest matching)
    - Network analysis (mutual connections)
    - Hybrid (combined scoring)
    """,
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# GZip compression
app.add_middleware(GZipMiddleware, minimum_size=1000)

# Include routers
app.include_router(recommendations.router, prefix="/api/v1")


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": "socialink-recommendation-service",
        "version": "1.0.0",
        "status": "operational",
        "recommendations": [
            "friends_suggestions",
            "people_you_may_know",
            "communities_for_you"
        ],
        "algorithms": [
            "collaborative_filtering",
            "graph_based",
            "content_based",
            "popularity",
            "hybrid"
        ]
    }


@app.get("/health")
async def health():
    """Health check"""
    return {
        "status": "healthy",
        "service": "socialink-recommendation-service",
        "version": "1.0.0"
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG
    )
