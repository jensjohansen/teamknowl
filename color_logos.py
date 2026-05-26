from PIL import Image, ImageOps
import os

def apply_color_overlay(input_path, output_path, color_hex, intensity=0.3):
    img = Image.open(input_path).convert("RGB")
    
    # Create a solid color layer
    color_img = Image.new("RGB", img.size, color_hex)
    
    # Blend the original image with the solid color
    tinted = Image.blend(img, color_img, intensity)
    
    # Create a mask from the original image to keep the background black
    # We'll treat anything very dark as background
    grayscale = img.convert("L")
    mask = grayscale.point(lambda x: 255 if x > 5 else 0, "L")
    
    # Composite: keep tinted only where original was not black
    final = Image.composite(tinted, img, mask)
    
    final.save(output_path)

assets_dir = "/home/johnk/.zenflow/worktrees/part-time-engineer/teamknowl/assets"
original = os.path.join(assets_dir, "knowl.png")

# Define target colors (hex)
colors = {
    "green": "#2ECC71",   # CodeKnowl
    "indigo": "#4B0082",  # TeamKnowl
    "orange": "#E67E22",  # DevOpsKnowl
    "blue": "#3498DB"     # SupportKnowl
}

intensity = 0.3

for name, color in colors.items():
    output = os.path.join(assets_dir, f"knowl-{name}.png")
    apply_color_overlay(original, output, color, intensity)
    print(f"Generated {output} with {color} at {intensity*100}% intensity")
