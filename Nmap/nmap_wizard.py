import json
import os

class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    END = '\033[0m'
    BOLD = '\033[1m'


OUTPUT_FILE = "selected_nmap_flags.json"

def save_to_json(data):

    try:
        # Ensure the data is in a valid format

        # Check if the output folder exists, create it if not
        output_folder = "intermediate_data"
        if not os.path.exists(output_folder):
            os.makedirs(output_folder)

        # Define the file path for the JSON file
        json_file_path = os.path.join(output_folder, OUTPUT_FILE)

        # Save the data to a JSON file with proper indentation
        with open(json_file_path, 'w') as json_file:
            json.dump(data, json_file, indent=2)

    except Exception as e:
        print(f"Error saving to JSON: {e}")

# Load flag data from JSON file
def load_flag_data(file_path="Nmap/flags_data.json"):
    if not os.path.exists(file_path):
        print(Colors.FAIL + f"Error: {file_path} does not exist." + Colors.END)
        exit(1)
    with open(file_path, "r") as file:
        return json.load(file)

# Handle value-required flags
def handle_flag_value(flag_details):
    if "value_prompt" in flag_details:
        print(Colors.CYAN + flag_details["value_prompt"] + Colors.END, end=' ')
        return input()
    else:
        print(Colors.FAIL + "Error: Value required but no prompt provided." + Colors.END)
        return None

# Commander Cody's category selection
def select_flags_commander(category):
    print(f"\n{Colors.BOLD}Category: {category['description']}{Colors.END}\n")
    if "greeting" in category:
        print(Colors.BLUE + category["greeting"] + Colors.END)

    # Print all available flags in this category
    print(Colors.WARNING + "Available flags:" + Colors.END)
    for i, (flag, details) in enumerate(category["flags"].items(), start=1):
        print(Colors.CYAN + str(i) +". " + Colors.END +flag +": " + details['description'])

    selected_flags = []
    flag_list = []  # Initialize flag_list
    choice = input(Colors.WARNING + "Enter the number(s) of your choice, separated by commas (or 'skip' to decline): " + Colors.END)

    if choice.lower() != 'skip':
        for i in choice.split(','):
            i = i.strip()
            if i.isdigit() and 1 <= int(i) <= len(category["flags"]):
                selected_flag = list(category["flags"].keys())[int(i)-1]
                flag_details = category["flags"][selected_flag]
                selected_flags.append({"flag": selected_flag})
                flag_list.append(selected_flag)  # Add to the list for consolidated printing
                if flag_details["value_required"]:
                    flag_value = handle_flag_value(flag_details)
                    selected_flags[-1]["value"] = flag_value
            else:
                print(Colors.FAIL + f"Invalid selection, trooper: {i}" + Colors.END)

    # Print all selected flags in a single line
    if flag_list:
        print(
            Colors.GREEN
            + f"Good choice, trooper. Flags {', '.join(flag_list)} are locked and loaded."
            + Colors.END
        )
    return selected_flags


# Main wizard function
def nmap_wizard_commander():
    print(Colors.GREEN + "Welcome, trooper! Commander Cody reporting. Let's prepare for the scan mission." + Colors.END)
    selected_flags = []
    flag_data = load_flag_data()

    for category_key, category_details in flag_data.items():
        category_flags = select_flags_commander(category_details)
        selected_flags.extend(category_flags)

    if selected_flags:
        # Output flags as JSON for Go to consume
        save_to_json(selected_flags)
    else:
        save_to_json("{}")


if __name__ == "__main__":
    nmap_wizard_commander()
