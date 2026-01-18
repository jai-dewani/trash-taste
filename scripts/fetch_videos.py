import os
import certifi
import ssl

import httplib2

# Patch httplib2 to use certifi certificates
original_init = httplib2.Http.__init__

def patched_init(self, *args, **kwargs):
    kwargs.setdefault('ca_certs', certifi.where())
    original_init(self, *args, **kwargs)

httplib2.Http.__init__ = patched_init

import json
from googleapiclient.discovery import build

from dotenv import load_dotenv
load_dotenv()
API_KEY = os.environ.get("YOUTUBE_DATA_API_V3")

# Trash Taste channel ID
CHANNEL_ID = "UCcmxOGYGF51T1XsqQLewGtQ"

def get_youtube_client():
    return build("youtube", "v3", developerKey=API_KEY)

def get_all_video_ids(youtube):
    video_ids = []
    
    # Get uploads playlist ID
    channel_response = youtube.channels().list(
        part="contentDetails",
        id=CHANNEL_ID
    ).execute()
    
    uploads_playlist_id = channel_response["items"][0]["contentDetails"]["relatedPlaylists"]["uploads"]
    
    # Fetch all videos from uploads playlist
    next_page_token = None
    while True:
        playlist_response = youtube.playlistItems().list(
            part="contentDetails",
            playlistId=uploads_playlist_id,
            maxResults=50,
            pageToken=next_page_token
        ).execute()
        
        for item in playlist_response["items"]:
            video_ids.append(item["contentDetails"]["videoId"])
        
        next_page_token = playlist_response.get("nextPageToken")
        if not next_page_token:
            break
    
    return video_ids

def get_video_details(youtube, video_ids):
    videos = []
    
    # Fetch details in batches of 50
    for i in range(0, len(video_ids), 50):
        batch_ids = video_ids[i:i+50]
        
        video_response = youtube.videos().list(
            part="snippet,contentDetails,statistics",
            id=",".join(batch_ids)
        ).execute()
        
        for item in video_response["items"]:
            video = {
                "id": item["id"],
                "title": item["snippet"]["title"],
                "description": item["snippet"]["description"],
                "publishedAt": item["snippet"]["publishedAt"],
                "thumbnails": item["snippet"]["thumbnails"],
                "tags": item["snippet"].get("tags", []),
                "duration": item["contentDetails"]["duration"],
                "viewCount": item["statistics"].get("viewCount", "0"),
                "likeCount": item["statistics"].get("likeCount", "0"),
                "commentCount": item["statistics"].get("commentCount", "0")
            }
            videos.append(video)
    
    return videos

def main():
    if not API_KEY:
        raise ValueError("YOUTUBE_DATA_API_V3 environment variable not set")
    
    youtube = get_youtube_client()
    
    print("Fetching video IDs...")
    video_ids = get_all_video_ids(youtube)
    print(f"Found {len(video_ids)} videos")
    
    print("Fetching video details...")
    videos = get_video_details(youtube, video_ids)
    print(f"Fetched details for {len(videos)} videos")
    
    # Ensure data directory exists
    os.makedirs("data", exist_ok=True)
    
    # Save to JSON file
    with open("data/videos.json", "w", encoding="utf-8") as f:
        json.dump(videos, f, indent=2, ensure_ascii=False)
    
    print("Saved videos to data/videos.json")

if __name__ == "__main__":
    main()