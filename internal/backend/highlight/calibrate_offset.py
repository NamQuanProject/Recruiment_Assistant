import sys
import os
import fitz
import json
from pathlib import Path

def create_calibration_pdf(pdf_path, output_dir, offset):
    """Create a calibration PDF with a specific y-offset."""
    # Open the PDF
    doc = fitz.open(pdf_path)
    
    # Create a new PDF for calibration
    output_path = os.path.join(output_dir, f"calibration_offset_{offset}.pdf")
    new_doc = fitz.open()
    
    # Process each page
    for page_num in range(len(doc)):
        page = doc[page_num]
        new_page = new_doc.new_page(width=page.rect.width, height=page.rect.height)
        
        # Copy the page content
        new_page.show_pdf_page(new_page.rect, doc, page_num)
        
        # Add a test highlight with the specified offset
        test_text = "Test Highlight"
        text_rect = fitz.Rect(50, 50 + offset, 200, 70 + offset)
        new_page.draw_rect(text_rect, color=(1, 1, 0), width=0.5)
        new_page.insert_text((50, 60 + offset), test_text, color=(0, 0, 0))
        
        # Add offset information
        info_text = f"Offset: {offset}"
        new_page.insert_text((50, 30), info_text, color=(0, 0, 0))
    
    # Save the calibration PDF
    new_doc.save(output_path)
    new_doc.close()
    doc.close()
    
    return output_path

def main():
    if len(sys.argv) != 3:
        print("Usage: python calibrate_offset.py <pdf_path> <output_dir>")
        sys.exit(1)
    
    pdf_path = sys.argv[1]
    output_dir = sys.argv[2]
    
    # Create output directory if it doesn't exist
    os.makedirs(output_dir, exist_ok=True)
    
    # Test different offsets
    offsets = [0, 1, 2, 3, 4, 5]
    calibration_files = []
    
    for offset in offsets:
        output_path = create_calibration_pdf(pdf_path, output_dir, offset)
        calibration_files.append({
            "offset": offset,
            "file": output_path
        })
    
    # Save calibration information
    calibration_info = {
        "pdf_path": pdf_path,
        "calibration_files": calibration_files
    }
    
    info_path = os.path.join(output_dir, "calibration_info.json")
    with open(info_path, "w") as f:
        json.dump(calibration_info, f, indent=2)
    
    print(f"Calibration completed. Check the following files:")
    print(f"Calibration info: {info_path}")
    for file_info in calibration_files:
        print(f"Offset {file_info['offset']}: {file_info['file']}")

if __name__ == "__main__":
    main() 