import subprocess

def main():
    result = subprocess.run(['make', 'init'], capture_output=True, text=True)
    print("Return code:", result.returncode)
    print("Output:", result.stdout)
    if result.stderr:
        print("Errors:", result.stderr)

if __name__ == "__main__":
    main()