"""
Vignette Feed Service
Main application entry point
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from contextlib import asynccontextmanager
import logging

from app.config import get_settings
from app.api.v1 import feed
from app.db import mongodb, redis_client


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
    logger.info("Starting Vignette Feed Service...")
    
    # Initialize database connections
    await mongodb.connect()
    await redis_client.connect()
    
    logger.info("✅ Database connections established")
    logger.info(f"🚀 Feed Service ready on port {settings.PORT}")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Feed Service...")
    await mongodb.disconnect()
    await redis_client.disconnect()
    logger.info("👋 Goodbye!")


# Create FastAPI app
app = FastAPI(
    title="Vignette Feed Service",
    description="""
    Advanced feed recommendation service with 3 feed algorithms:
    
    - **Home Feed**: TikTok-style For You feed (viral discovery)
    - **Circle Feed**: Friends + nearby content (social graph)
    - **Surprise & Delight**: Balanced exploration (60% interests, 30% explore, 10% surprise)
    
    Features:
    - Personalized ranking
    - Real-time signal tracking
    - Interest extraction
    - Collaborative filtering
    - Location-based recommendations
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
app.include_router(feed.router, prefix="/api/v1")


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": "vignette-feed-service",
        "version": "1.0.0",
        "status": "operational",
        "feeds": [
            "home",
            "circle",
            "surprise_delight"
        ],
        "description": "Advanced feed recommendation with 3 algorithms"
    }


@app.get("/health")
async def health():
    """Health check"""
    return {
        "status": "healthy",
        "service": "vignette-feed-service",
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
