import fitz  # PyMuPDF
import argparse
import os
import sys

def parse_pdf(pdf_path, output_path):
    doc = fitz.open(pdf_path)
    full_text = ""

    for page_num in range(len(doc)):
        page = doc[page_num]
        words = page.get_text("words")  # Get all words with positions
        links = page.get_links()        # Get all link annotations

        word_objs = [{
            'text': w[4],
            'rect': fitz.Rect(w[0], w[1], w[2], w[3]),
            'used': False
        } for w in words]

        for link in links:
            if 'uri' not in link:
                continue
            link_rect = fitz.Rect(link['from'])
            link_text_indices = []

            for i, w in enumerate(word_objs):
                if not w['used'] and w['rect'].intersects(link_rect):
                    link_text_indices.append(i)
                    w['used'] = True

            if link_text_indices:
                last_index = link_text_indices[-1]
                word_objs[last_index]['text'] += f" ({link['uri']})"

        sorted_words = sorted(word_objs, key=lambda w: (round(w['rect'].y0, 1), w['rect'].x0))
        text = ""
        prev_y = None

        for w in sorted_words:
            y = round(w['rect'].y0, 1)
            if prev_y is not None and y != prev_y:
                text += "\n"
            elif text and not text.endswith(" "):
                text += " "
            text += w['text']
            prev_y = y

        full_text += f"\n--- Page {page_num + 1} ---\n{text}\n"

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(full_text)

    print(f"Text with links saved to: {output_path}")


def parse_pdf_batch(input_dir, output_dir):
    os.makedirs(output_dir, exist_ok=True)

    for file_name in os.listdir(input_dir):
        if file_name.lower().endswith(".pdf"):
            input_path = os.path.join(input_dir, file_name)
            output_name = os.path.splitext(file_name)[0] + ".txt"
            output_path = os.path.join(output_dir, output_name)

            print(f"Processing {file_name}...")
            parse_pdf(input_path, output_path)

    print(f"All PDFs in {input_dir} have been converted to TXTs in {output_dir}.")


def main():
    parser = argparse.ArgumentParser(description="Extract text and hyperlinks from PDF(s).")

    parser.add_argument("-batch", help="Enable batch mode", default="false")
    parser.add_argument("input", help="PDF file path or input directory")
    parser.add_argument("output", help="TXT file path or output directory")

    args = parser.parse_args()
    batch_mode = args.batch.lower() == "true"

    if batch_mode:
        if not os.path.isdir(args.input):
            print(f"Error: Batch mode expects input directory, got file: {args.input}")
            sys.exit(1)
        parse_pdf_batch(args.input, args.output)
    else:
        if not os.path.isfile(args.input):
            print(f"Error: File not found: {args.input}")
            sys.exit(1)
        parse_pdf(args.input, args.output)

if __name__ == "__main__":
    main()
