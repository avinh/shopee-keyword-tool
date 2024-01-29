const form = document.getElementById('control-row');
const go = document.getElementById('go');

(async function initPopupWindow() {
})();

form.addEventListener('submit', handleFormSubmit);

async function handleFormSubmit(event) {
  event.preventDefault();

  let url = new URL('https://banhang.shopee.vn');
  if (!url) {
    return;
  }

  let message = await getDomainCookies(url.hostname);
  prompt('Sao chép vào clipboard: Ctrl+C, Enter', message);
}

async function getDomainCookies(domain) {
  try {
    const cookies = await chrome.cookies.getAll({ domain });
    if (cookies.length === 0) {
      return 'No cookies found';
    }
    const cookieString = cookies.map(cookie => `${cookie.name}=${cookie.value}`).join('; ');
    return cookieString;
  } catch (error) {
    return `Unexpected error: ${error.message}`;
  }
}

function getCookieInfo(cookie) {
  return JSON.stringify({ name: cookie.name, value: cookie.value });
}
