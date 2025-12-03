import datetime
import os

def main():
    # Get the current date and time
    now = datetime.datetime.utcnow().strftime("%Y%m%d%H%M%S")
    
    # Define the file where you want to insert the timestamp
    target_file = "cookiecutter.json"  # Change this to your desired file path
    
    # Check if the file exists
    if os.path.exists(target_file):
        # Read the file content
        with open(target_file, 'r') as file:
            content = file.read()
        
        # Replace a placeholder with the current timestamp
        # Example: Replace '{{CURRENT_DATE}}' with the actual date/time
        updated_content = content.replace('{{ now }}', now)
        
        # Write the updated content back to the file
        with open(target_file, 'w') as file:
            file.write(updated_content)
    else:
        print(f"Warning: {target_file} not found.")

if __name__ == "__main__":
    main()