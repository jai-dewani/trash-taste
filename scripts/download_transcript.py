import os
import json
from youtube_transcript_api import YouTubeTranscriptApi

def transcript_exists(video_id, output_dir):
    """Check if transcript file already exists."""
    output_path = os.path.join(output_dir, f"{video_id}.json")
    return os.path.exists(output_path)

def load_videos(file_path):
    """Load list of video IDs from a file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def download_transcript(video_id, output_dir):
    """Download transcript for a video and save to file."""
    if transcript_exists(video_id, output_dir):
        print(f"Transcript for {video_id} already exists. Skipping download.")
        return True
    try:
        transcript = YouTubeTranscriptApi().fetch(video_id)
        
        # Convert FetchedTranscript to a list of dictionaries
        transcript_data = [
            {
                'text': snippet.text,
                'start': snippet.start,
                'duration': snippet.duration
            }
            for snippet in transcript
        ]
        
        output_path = os.path.join(output_dir, f"{video_id}.json")
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(transcript_data, f, ensure_ascii=False, indent=2)
        
        print(f"Downloaded transcript for {video_id}")
        return True
    except Exception as e:
        print(f"Failed to download transcript for {video_id}: {e}")
        return False

def main():
    videos_file = "../data/videos.json"
    output_dir = "../data/transcripts"
    
    os.makedirs(output_dir, exist_ok=True)
    
    videos = load_videos(videos_file)
    
    success_count = 0
    fail_count = 0
    
    for video in videos:
        video_id = video if isinstance(video, str) else video.get('id')
        if video_id:
            if download_transcript(video_id, output_dir):
                success_count += 1
            else:
                fail_count += 1
    
    print(f"\nCompleted: {success_count} successful, {fail_count} failed")

if __name__ == "__main__":
    main()