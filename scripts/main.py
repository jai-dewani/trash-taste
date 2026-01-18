import os
from dotenv import load_dotenv

# Load the environment variables from .env file
load_dotenv()

# Now you can access the environment variable just like before
api_key = os.environ.get('YOUTUBE_DATA_API_V3')
print(f"Your YouTube Data API v3 Key is: {api_key}")


