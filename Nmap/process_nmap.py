import os
import json
import xml.etree.ElementTree as ET


def parse_nmap_xml(xml_file):
    """
    Parses an Nmap XML file and converts the entire XML tree into a JSON-compatible dictionary.
    """
    def xml_to_dict(element):
        """
        Recursively converts an XML element and its children into a dictionary.
        """
        node = {}
        # Add element attributes
        if element.attrib:
            node.update({"@attributes": element.attrib})

        # Add element text content
        if element.text and element.text.strip():
            node.update({"@text": element.text.strip()})

        # Recursively process child elements
        children = list(element)
        if children:
            child_dict = {}
            for child in children:
                child_name = child.tag
                child_data = xml_to_dict(child)

                # Handle multiple children with the same tag
                if child_name in child_dict:
                    if not isinstance(child_dict[child_name], list):
                        child_dict[child_name] = [child_dict[child_name]]
                    child_dict[child_name].append(child_data)
                else:
                    child_dict[child_name] = child_data

            node.update({"@children": child_dict})

        return node

    # Parse the XML file
    try:
        tree = ET.parse(xml_file)
        root = tree.getroot()
        return json.dumps(xml_to_dict(root), indent=4)
    except ET.ParseError as e:
        raise ValueError(f"Error parsing XML file: {e}")
    except Exception as e:
        raise ValueError(f"Unexpected error during XML parsing: {e}")


def get_filenames_from_record():
    """
    Retrieve file paths from the filenames record JSON.
    """
    filenames_record_file = os.path.join(os.path.dirname(__file__), '..', 'filenames.json')

    if not os.path.exists(filenames_record_file):
        raise FileNotFoundError("The filenames record 'filenames.json' does not exist.")

    with open(filenames_record_file, 'r') as f:
        try:
            filenames_record = json.load(f)
        except json.JSONDecodeError:
            raise ValueError("The filenames record is not a valid JSON file.")

    for record in filenames_record:
        if 'Nmap' in record:
            return record['Nmap']

    raise ValueError("No 'Nmap' entry found in the filenames record.")

def save_to_json_file(parsed_result, nmap_files):
    """
    Save the parsed JSON data to a file in the output folder.
    """
    output_folder = "output"
    os.makedirs(output_folder, exist_ok=True)

    # Construct the JSON file path
    json_file_path = os.path.join(output_folder, nmap_files['json'])

    # Avoid overwriting an existing JSON file
    if os.path.exists(json_file_path):
        raise FileExistsError(f"The JSON file '{nmap_files['json']}' already exists. Remove it or change the filename.")

    # Save the parsed result to the JSON file
    with open(json_file_path, 'w') as json_file:
        json_file.write(parsed_result)


if __name__ == "__main__":
    try:
        nmap_xml_path = "output/nmap_xml.xml"  # Path to the existing XML file
        nmap_files = get_filenames_from_record()  # Get filenames from the record

        # Parse the Nmap XML file and convert it to JSON
        parsed_json = parse_nmap_xml(nmap_xml_path)

        # Save only the JSON file
        save_to_json_file(parsed_json, nmap_files)

        print("Nmap JSON file saved successfully.")
    except FileNotFoundError as e:
        print(f"File error: {e}")
    except ValueError as e:
        print(f"Data error: {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
