"""
Vignette Recommendation Service
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
    logger.info("Starting Vignette Recommendation Service...")
    
    # Initialize ML models, DB connections, etc.
    # await initialize_services()
    
    logger.info(f"✅ Recommendation Service ready on port {settings.PORT}")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Recommendation Service...")
    logger.info("👋 Goodbye!")


# Create FastAPI app
app = FastAPI(
    title="Vignette Recommendation Service",
    description="""
    ML-powered recommendation system for Vignette
    
    Provides 3 types of recommendations:
    - **Creators for You**: Top creators to follow
    - **Suggested for You**: General user recommendations
    - **Communities for You**: Community recommendations
    
    Uses multiple algorithms:
    - Collaborative filtering (similar users)
    - Graph-based (friends of friends, PageRank)
    - Content-based (interest matching)
    - Popularity-based (trending)
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
        "service": "vignette-recommendation-service",
        "version": "1.0.0",
        "status": "operational",
        "recommendations": [
            "creators_for_you",
            "suggested_for_you",
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
        "service": "vignette-recommendation-service",
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
