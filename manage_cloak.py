import argparse
import base64
import json
import requests
import os
import sys
import uuid

def set_api_base(base_url=None):
    api_base = base_url or os.getenv('API_BASE')
    if not api_base:
        print("Error: API base URL is not set. Provide it via --base-url or the API_BASE environment variable.")
        sys.exit(1)
    return api_base

def alert_error(response):
    print(f"Error: {response.status_code} {response.text}")
    sys.exit(1)

def base64_url_encode(data):
    return base64.urlsafe_b64encode(data.encode()).decode().replace('=', '')

def generate_uid():
    return str(uuid.uuid4())

def list_users(api_base):
    endpoint = f"http://{api_base}/admin/users"
    try:
        response = requests.get(endpoint)
        response.raise_for_status()
        users = response.json()
        for user in users:
            print(json.dumps(user, indent=4))
    except requests.RequestException as e:
        alert_error(e.response)

def delete_user(api_base, uid):
    encoded_uid = base64_url_encode(uid)
    endpoint = f"http://{api_base}/admin/users/{encoded_uid}"
    confirm = input(f"Do you really want to delete {uid}? (yes/no): ").strip().lower()
    if confirm != "yes":
        print("User deletion canceled.")
        return
    try:
        response = requests.delete(endpoint)
        response.raise_for_status()
        print(f"User {uid} deleted successfully.")
    except requests.RequestException as e:
        alert_error(e.response)

def add_user(api_base, userinfo):
    encoded_uid = base64_url_encode(userinfo['UID'])
    endpoint = f"http://{api_base}/admin/users/{encoded_uid}"
    try:
        response = requests.post(endpoint, json=userinfo)
        response.raise_for_status()
        print("User added successfully.")
    except requests.RequestException as e:
        alert_error(e.response)

def main():
    parser = argparse.ArgumentParser(description="CLI tool to manage users via API.")
    parser.add_argument('--base-url', type=str, help="Base URL for the API (or set API_BASE environment variable).")
    subparsers = parser.add_subparsers(dest="command", required=True)

    subparsers.add_parser('list-users', help="List all users.")

    delete_parser = subparsers.add_parser('delete-user', help="Delete a user.")
    delete_parser.add_argument('uid', type=str, help="UID of the user to delete.")

    add_parser = subparsers.add_parser('add-user', help="Add a new user.")
    add_parser.add_argument('--uid', type=str, help="UID for the user. Generated if not provided.")
    add_parser.add_argument('--sessions-cap', type=int, required=True, help="Session cap for the user.")
    add_parser.add_argument('--up-rate', type=int, required=True, help="Upload rate in Mbps.")
    add_parser.add_argument('--down-rate', type=int, required=True, help="Download rate in Mbps.")
    add_parser.add_argument('--up-credit', type=int, required=True, help="Upload credit in MB.")
    add_parser.add_argument('--down-credit', type=int, required=True, help="Download credit in MB.")
    add_parser.add_argument('--expiry-time', type=int, required=True, help="Expiry time in seconds.")

    args = parser.parse_args()
    api_base = set_api_base(args.base_url)

    if args.command == "list-users":
        list_users(api_base)
    elif args.command == "delete-user":
        delete_user(api_base, args.uid)
    elif args.command == "add-user":
        userinfo = {
            "UID": args.uid or generate_uid(),
            "SessionsCap": args.sessions_cap,
            "UpRate": args.up_rate * 1048576,  # Convert Mbps to bits
            "DownRate": args.down_rate * 1048576,  # Convert Mbps to bits
            "UpCredit": args.up_credit * 1048576,  # Convert MB to bits
            "DownCredit": args.down_credit * 1048576,  # Convert MB to bits
            "ExpiryTime": args.expiry_time,
        }
        add_user(api_base, userinfo)

if __name__ == "__main__":
    main()
