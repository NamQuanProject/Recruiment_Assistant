import fitz  # PyMuPDF
import json
import argparse
import sys

def highlight_pdf(pdf_path, areas_path, output_path):
    """
    Highlight strong and weak areas in a PDF file.
    
    Args:
        pdf_path (str): Path to the input PDF file
        areas_path (str): Path to the JSON file containing strong and weak areas
        output_path (str): Path to save the highlighted PDF
    """
    # Load the PDF
    doc = fitz.open(pdf_path)
    
    # Load areas from JSON
    with open(areas_path, 'r', encoding='utf-8') as f:
        areas = json.load(f)
    
    # Highlight each area
    for area in areas:
        page_num = area.get('page', 0) - 1  # Convert to 0-based index
        
        # Skip if page number is out of range
        if page_num < 0 or page_num >= len(doc):
            print(f"Warning: Page {area.get('page')} is out of range. Skipping.")
            continue
        
        page = doc[page_num]
        
        # Apply a small offset to the y-coordinate to account for PDF rendering differences
        y_offset = 0  # Adjust this value if needed (positive values move highlight down)
        
        # Create a highlight annotation with precise coordinates
        rect = fitz.Rect(
            area['x'], 
            area['y'] + y_offset, 
            area['x'] + area['width'], 
            area['y'] + area['height'] + y_offset
        )
        
        # Add a highlight annotation with color based on type
        annot = page.add_highlight_annot(rect)
        
        # Set color based on type (strong = green, weak = yellow)
        area_type = area.get('type', 'weak').lower()
        if area_type == 'strong':
            annot.set_colors(stroke=[0, 1, 0])  # Green color
            title = "Strong Area"
        else:
            annot.set_colors(stroke=[1, 1, 0])  # Yellow color
            title = "Weak Area"
            
        annot.set_opacity(0.3)  # 30% opacity
        
        # Add a popup note with the description
        if 'description' in area and area['description']:
            annot.set_info(title=title, content=area['description'])
        
        # Update the annotation
        annot.update()
    
    # Save the highlighted PDF
    doc.save(output_path)
    doc.close()
    
    print(f"Highlighted PDF saved to: {output_path}")

def main():
    parser = argparse.ArgumentParser(description="Highlight strong and weak areas in a PDF file.")
    parser.add_argument("pdf_path", help="Path to the input PDF file")
    parser.add_argument("areas_path", help="Path to the JSON file containing strong and weak areas")
    parser.add_argument("output_path", help="Path to save the highlighted PDF")
    
    args = parser.parse_args()
    
    highlight_pdf(args.pdf_path, args.areas_path, args.output_path)

if __name__ == "__main__":
    main() 