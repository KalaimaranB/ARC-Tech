import sys
import re
import json
import os


def parse_nmap_output(nmap_output):
    # Initialize the result dictionary
    result = {
        'ports': [],
        'service_info': '',
        'scan_summary': ''
    }

    # Regular expressions to capture port details and service info
    port_pattern = re.compile(r"(\d+/tcp)\s+(\S+)\s+(\S+)\s*(.*)")
    service_info_pattern = re.compile(r"Service Info: (.*)")
    scan_summary_pattern = re.compile(r"Nmap done:.*scanned in (.*) seconds")

    # Split the input into lines
    lines = nmap_output.splitlines()

    # Process each line
    for line in lines:
        # Match port details
        port_match = port_pattern.match(line)
        if port_match:
            port = port_match.group(1)
            state = port_match.group(2)
            service = port_match.group(3)
            version = port_match.group(4) if port_match.group(4) else ""  # Use empty string if version is missing
            result['ports'].append({
                'port': port,
                'state': state,
                'service': service,
                'version': version
            })

        # Match service info
        service_info_match = service_info_pattern.match(line)
        if service_info_match:
            result['service_info'] = service_info_match.group(1)

        # Match scan summary
        scan_summary_match = scan_summary_pattern.search(line)
        if scan_summary_match:
            result['scan_summary'] = scan_summary_match.group(0)

    # Return the parsed result as a JSON string
    return json.dumps(result, indent=4)


def get_filenames_from_record():
    # Define the filename for the filenames record
    filenames_record_file = os.path.join(os.path.dirname(__file__), '..', 'filenames.json')

    # If the file doesn't exist, create it with an empty list
    if not os.path.exists(filenames_record_file):
        filenames_record = []
    else:
        # Read existing filenames from the file
        with open(filenames_record_file, 'r') as f:
            filenames_record = json.load(f)

    # Find the "Nmap" entry in the record
    for record in filenames_record:
        if 'Nmap' in record:
            return record['Nmap']

    # If no Nmap entry exists, raise an error
    raise ValueError("No Nmap entry found in filenames record.")


def save_to_files(nmap_output, parsed_result, nmap_files):
    # Define the output folder (same level as the script)
    output_folder = "output"

    # Check if the output folder exists, create it if not
    if not os.path.exists(output_folder):
        os.makedirs(output_folder)

    # Construct paths for the text and json files
    text_file_path = os.path.join(output_folder, nmap_files['text'])
    json_file_path = os.path.join(output_folder, nmap_files['json'])

    # Save the raw Nmap output to a .txt file
    with open(text_file_path, 'w') as text_file:
        text_file.write(nmap_output)

    # Save the parsed result to the JSON file
    with open(json_file_path, 'w') as json_file:
        json_file.write(parsed_result)


def main():
    print("Starting python analysis")

    # Read the entire content from stdin (fix for your error)
    nmap_output = sys.stdin.read()

    # Parse the Nmap output
    parsed_result = parse_nmap_output(nmap_output)

    # Get the file names from the filenames record
    nmap_files = get_filenames_from_record()

    # Save both the raw Nmap output and the parsed result to files
    save_to_files(nmap_output, parsed_result, nmap_files)
    print("Finished python analysis")


if __name__ == "__main__":
    main()
