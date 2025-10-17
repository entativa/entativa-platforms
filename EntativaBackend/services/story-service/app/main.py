"""
Socialink Story Service - Main Application
Ephemeral stories with interactive stickers, highlights, and analytics
"""
import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn

from app.config import settings
from app.db.mongodb import MongoDB
from app.db.redis_client import RedisClient
from app.services.expiration_service import ExpirationService
from app.api import stories_router, highlights_router, viewers_router

# Setup logging
logging.basicConfig(
    level=logging.DEBUG if settings.DEBUG else logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

# Expiration service instance
expiration_service = ExpirationService()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan events for startup and shutdown"""
    # Startup
    logger.info(f"Starting {settings.APP_NAME} v{settings.APP_VERSION}")
    
    # Connect to databases
    await MongoDB.connect()
    await RedisClient.connect()
    
    # Start expiration service in background
    asyncio.create_task(expiration_service.start())
    
    logger.info("Story service started successfully")
    
    yield
    
    # Shutdown
    logger.info("Shutting down story service...")
    
    expiration_service.stop()
    await MongoDB.close()
    await RedisClient.close()
    
    logger.info("Story service shut down successfully")


# Create FastAPI app
app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    description="""
    **Socialink Story Service** 🎬
    
    Comprehensive story management with:
    - ⏱️ 24-hour ephemeral stories
    - 🎯 Interactive stickers (polls, quizzes, questions, countdowns, sliders)
    - 📌 Permanent highlights
    - 👥 Close friends feature
    - 📊 Analytics and insights
    - 💬 Story replies
    - 👀 View tracking
    
    ## Features
    
    ### Stories
    - Create image, video, or text stories
    - Add interactive stickers for engagement
    - Privacy controls (public, followers, close friends)
    - Auto-expiration after 24 hours
    - View tracking with "seen by" lists
    
    ### Interactive Stickers
    - **Polls**: Multiple choice voting with percentages
    - **Quizzes**: Test followers with right/wrong answers
    - **Questions**: Open-ended responses
    - **Countdowns**: Events with follower notifications
    - **Sliders**: Emoji sliders for ratings
    - **Mentions**: Tag other users
    - **Location**: Share where you are
    - **Music**: Add soundtrack
    
    ### Highlights
    - Create permanent story collections
    - Custom covers and titles
    - Organize and reorder
    - Pin important highlights
    
    ### Analytics
    - View counts and reach
    - Interaction rates
    - Peak viewing times
    - Sticker performance
    - Top viewers
    """,
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure properly in production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# Exception handler
@app.exception_handler(Exception)
async def global_exception_handler(request: Request, exc: Exception):
    """Global exception handler"""
    logger.error(f"Global exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=500,
        content={"detail": "Internal server error"}
    )


# Health check
@app.get("/health", tags=["health"])
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": settings.APP_NAME,
        "version": settings.APP_VERSION
    }


# Root endpoint
@app.get("/", tags=["root"])
async def root():
    """Root endpoint"""
    return {
        "service": settings.APP_NAME,
        "version": settings.APP_VERSION,
        "description": "Socialink Story Service - Ephemeral stories with interactive features",
        "features": [
            "24-hour stories",
            "Interactive stickers (polls, quizzes, questions, etc.)",
            "Story highlights",
            "Close friends",
            "Analytics",
            "View tracking"
        ],
        "docs": "/docs"
    }


# Register routers
app.include_router(stories_router, prefix="/api/v1")
app.include_router(highlights_router, prefix="/api/v1")
app.include_router(viewers_router, prefix="/api/v1")

# Close friends endpoints (inline for simplicity)
from fastapi import APIRouter, Header
from app.services.close_friends_service import CloseFriendsService

close_friends_router = APIRouter(prefix="/api/v1/close-friends", tags=["close-friends"])
close_friends_service = CloseFriendsService()


@close_friends_router.post("/{friend_user_id}")
async def add_close_friend(
    friend_user_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Add a user to close friends list"""
    result = await close_friends_service.add_close_friend(x_user_id, friend_user_id)
    return {"success": result, "message": "Added to close friends" if result else "Already in close friends"}


@close_friends_router.delete("/{friend_user_id}")
async def remove_close_friend(
    friend_user_id: str,
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Remove a user from close friends list"""
    result = await close_friends_service.remove_close_friend(x_user_id, friend_user_id)
    return {"success": result, "message": "Removed from close friends" if result else "Not in close friends"}


@close_friends_router.get("/")
async def get_close_friends(
    x_user_id: str = Header(..., description="User ID from auth")
):
    """Get close friends list"""
    friends = await close_friends_service.get_close_friends(x_user_id)
    return {"close_friends": friends, "count": len(friends)}


app.include_router(close_friends_router)


if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="debug" if settings.DEBUG else "info"
    )
