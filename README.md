# Shopee Keyword Recommender

This is a Go CLI tool that helps Shopee sellers retrieve keyword recommendations and their suggested prices for a selected product. It interacts with Shopee's internal APIs to get keyword hints and keyword data.

## Features

- Automatically fetches a product list from Shopee's seller portal.
- Retrieves keyword suggestions based on an input keyword and selected product.
- Fetches detailed keyword information, including:
  - Recommended price
  - Relevance
  - Search volume
- Saves the output in `output.csv`.

---

## Setup

### 1. Install Chrome Extension to Get Shopee Cookie

Inside the `chrome-extensions` folder, there's a Chrome extension to help you extract your Shopee seller account cookie.

#### Steps:

1. Open Google Chrome.
2. Go to `chrome://extensions/`.
3. Enable "Developer mode" (top right corner).
4. Click "Load unpacked" and select the `chrome-extensions` folder.
5. After installing the extension, log in to your Shopee seller account.
6. Click the extension icon to copy your session cookie.

### 2. Configure `config.json`

Create or edit the `config.json` file in the root directory with the following format:

```json
{
  "cookie": "YOUR_SHOPEE_COOKIE_HERE"
}
```

Replace `YOUR_SHOPEE_COOKIE_HERE` with the cookie value obtained using the Chrome extension.

---

## Usage

Build and run the program:

```bash
go run main.go
```

### Example Flow

1. The program will ask:  
   `Please enter keyword:`  
   You type: `shoes`

2. It fetches keyword hints and data from Shopee for the first product.

3. Keyword data (keyword and recommended price) will be saved in `output.csv`.

---

## Output

The `output.csv` file will have the following structure:

```
keyword,recommended_price
running shoes,1500
sports shoes,1800
...
```

---

## Notes

- This tool accesses internal Shopee APIs. Use responsibly and only for your own seller account.
- Make sure your Shopee session is valid. If you encounter issues, try updating your cookie again.
