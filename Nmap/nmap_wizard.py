import json
import sys

# Color codes for terminal styling
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def load_flag_data(file_path='Nmap/flags_data.json'):
    """Load flag categories and descriptions from a JSON file."""
    try:
        with open(file_path, 'r') as file:
            return json.load(file)
    except FileNotFoundError:
        print(f"{Colors.FAIL}Error: The file {file_path} was not found.{Colors.ENDC}")
        sys.exit(1)
    except json.JSONDecodeError:
        print(f"{Colors.FAIL}Error: Failed to parse JSON data.{Colors.ENDC}")
        sys.exit(1)

def ask_question(question, options):
    """Generic function to ask a question and return the user's choice."""
    print(Colors.HEADER + question + Colors.ENDC)
    for i, (key, value) in enumerate(options.items(), start=1):
        print(f"{i}. {Colors.OKCYAN}{key}{Colors.ENDC}: {value}")
    print(Colors.WARNING + "Enter the number(s) of your choice, separated by commas (or 'skip' to decline): " + Colors.ENDC, end='')
    choice = input()
    return choice

def select_flags(category):
    """Select specific flags based on user input."""
    print(f"\n{Colors.BOLD}Category: {category['description']}{Colors.ENDC}\n")
    selected_flags = []

    # List out all flags and their descriptions
    for i, (flag, details) in enumerate(category["flags"].items(), start=1):
        print(f"{i}. {Colors.OKCYAN}{flag}{Colors.ENDC}: {details['description']}")

    # Ask user for selection
    choice = input(Colors.WARNING + "Enter the number(s) of your choice, separated by commas (or 'skip' to decline): " + Colors.ENDC)
    if choice.lower() != 'skip':
        for i in choice.split(','):
            i = i.strip()
            if i.isdigit() and 1 <= int(i) <= len(category["flags"]):
                selected_flag = list(category["flags"].keys())[int(i)-1]
                flag_details = category["flags"][selected_flag]
                selected_flags.append({"flag": selected_flag, "details": flag_details})
            else:
                print(Colors.FAIL + f"Invalid selection: {i}" + Colors.ENDC)

    return selected_flags

def handle_flag_value(flag_details):
    """If the flag requires a value, ask the user for it."""
    if flag_details["value_required"]:
        print(Colors.HEADER + flag_details["value_prompt"] + Colors.ENDC)
        value = input(Colors.WARNING + "Enter value: " + Colors.ENDC)
        return value
    return None

def nmap_wizard():
    """Main wizard function to guide the user through selecting flags."""
    selected_flags = []

    # Load data from JSON
    flag_data = load_flag_data()

    # Ask about network evasion preferences
    network_evasion_choice = ask_question("Do you care about avoiding detection during your scan?",
                                          {"1": "Network evasion flags help avoid detection.", "2": "No concern for evasion."})
    if network_evasion_choice == '1':
        selected_flags.extend(select_flags(flag_data["network_evasion"]))

    # Ask about packet manipulation preferences
    packet_manipulation_choice = ask_question("Do you want to manipulate packets to avoid detection?",
                                              {"1": "Packet manipulation flags help modify packet behavior.", "2": "No concern for packet manipulation."})
    if packet_manipulation_choice == '1':
        selected_flags.extend(select_flags(flag_data["packet_manipulation"]))

    # Ask about speed preferences
    speed_choice = ask_question("Do you care about scan speed?", {"1": "Speed flags optimize scan times.", "2": "Scan speed isn't a priority."})
    if speed_choice == '1':
        selected_flags.extend(select_flags(flag_data["speed"]))

    # Ask about detail preferences
    detailed_choice = ask_question("Do you need detailed information?", {"1": "Detailed flags provide in-depth analysis of services.", "2": "You prefer a more basic scan."})
    if detailed_choice == '1':
        selected_flags.extend(select_flags(flag_data["detailed"]))

    # Ask about aggressive scan preferences
    aggressive_choice = ask_question("Do you want to run an aggressive scan?",
                                     {"1": "Aggressive flags maximize scan results but may be loud.", "2": "Prefer a quieter scan."})
    if aggressive_choice == '1':
        selected_flags.extend(select_flags(flag_data["aggressive"]))

    # Ask for values if necessary
    final_flags = []
    for selected_flag in selected_flags:
        flag_value = handle_flag_value(selected_flag["details"])
        if flag_value:
            final_flags.append(f"{selected_flag['flag']}={flag_value}")
        else:
            final_flags.append(selected_flag['flag'])

    # Output selected flags with values
    print("\n" + Colors.OKGREEN + Colors.BOLD + "Selected Flags:" + Colors.ENDC)
    if final_flags:
        for flag in final_flags:
            print(f"{Colors.OKCYAN}{flag}{Colors.ENDC}")
    else:
        print(Colors.FAIL + "No flags selected." + Colors.ENDC)


if __name__ == "__main__":
    nmap_wizard()
