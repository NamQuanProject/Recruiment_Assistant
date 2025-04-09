import fitz  # PyMuPDF
import json
import argparse
import sys

def extract_text_from_pdf(pdf_path, output_path):
    """
    Extract text and positions from a PDF file.
    
    Args:
        pdf_path (str): Path to the input PDF file
        output_path (str): Path to save the extracted text blocks as JSON
    """
    # Load the PDF
    doc = fitz.open(pdf_path)
    
    # Extract text blocks from each page
    text_blocks = []
    
    for page_num, page in enumerate(doc):
        # Get text blocks with their positions
        blocks = page.get_text("dict")["blocks"]
        
        for block in blocks:
            # Skip non-text blocks
            if "lines" not in block:
                continue
                
            # Process each line in the block
            for line in block["lines"]:
                for span in line["spans"]:
                    # Get the exact bounding box
                    bbox = span["bbox"]
                    
                    # Create a text block with more precise coordinates
                    text_block = {
                        "text": span["text"].strip(),
                        "page": page_num + 1,  # Convert to 1-based index
                        "x": bbox[0],  # Left edge
                        "y": bbox[1],  # Top edge
                        "width": bbox[2] - bbox[0],  # Width
                        "height": bbox[3] - bbox[1]  # Height
                    }
                    
                    # Only add non-empty text blocks
                    if text_block["text"]:
                        text_blocks.append(text_block)
    
    # Save the text blocks to a JSON file
    with open(output_path, 'w') as f:
        json.dump(text_blocks, f, indent=2)
    
    print(f"Extracted {len(text_blocks)} text blocks from {pdf_path}")
    print(f"Text blocks saved to {output_path}")

def main():
    parser = argparse.ArgumentParser(description="Extract text and positions from a PDF file.")
    parser.add_argument("pdf_path", help="Path to the input PDF file")
    parser.add_argument("output_path", help="Path to save the extracted text blocks as JSON")
    
    args = parser.parse_args()
    
    extract_text_from_pdf(args.pdf_path, args.output_path)

if __name__ == "__main__":
    main() 