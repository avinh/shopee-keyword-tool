const form = document.getElementById('control-row');
const go = document.getElementById('go');
const cookieResult = document.getElementById('cookie');

(async function initPopupWindow() {
})();

form.addEventListener('submit', handleFormSubmit);

async function handleFormSubmit(event) {
  event.preventDefault();

  let url = new URL('https://banhang.shopee.vn');

  let message = await getDomainCookies(url.origin);

  // copy to clipboard
  const input = document.createElement('input');
  input.setAttribute('value', message);
  document.body.appendChild(input);
  input.select();
  const result = document.execCommand('copy');
  document.body.removeChild(input);

  if (result) {
    go.innerText = 'Đã copy vào clipboard!';
  }

  setTimeout(() => {
    go.innerText = 'Get cookie';
  }, 3000);

  cookieResult.innerText = message;
}

function Filter() {
  var filter = {};

  this.setUrl = function (url) {
    filter.url = url;
  };

  this.setDomain = function (domain) {
    filter.domain = domain;
  };
  this.setName = function (name) {
    filter.name = name;
  };
  this.setSecure = function (secure) {
    filter.secure = secure;
  };
  this.setSession = function (session) {
    filter.session = session;
  };
  this.getFilter = function (session) {
    return filter;
  };
}

async function getDomainCookies(url) {
  try {

    var filter = new Filter();
    if (/^https?:\/\/.+$/.test(url)) {
      filter.setUrl(url);
    } else {
      filter.setDomain(url);
    }

    const filters = filter.getFilter();

    if (filters === null)
      filters = {};
    var filterURL = {};
    if (filters.url !== undefined)
      filterURL.url = filters.url;
    if (filters.domain !== undefined)
      filterURL.domain = filters.domain;

    const cookieStores = await new Promise((resolve, reject) => {
      try {
        chrome.cookies.getAllCookieStores(function (cookieStores) {
          resolve(cookieStores);
        });
      } catch (error) {
        reject(error);
      }
    });

    if (cookieStores.length === 0)
      return 'No cookie store found';

    filterURL.storeId = cookieStores[0].id;

    const list = await new Promise((resolve, reject) => {
      try {
        const filteredCookies = [];
        chrome.cookies.getAll(filterURL, function (cks) {
          var currentC;
          for (var i = 0; i < cks.length; i++) {
            currentC = cks[i];

            if (filters.name !== undefined && currentC.name.toLowerCase().indexOf(filters.name.toLowerCase()) === -1)
              continue;
            if (filters.domain !== undefined && currentC.domain.toLowerCase().indexOf(filters.domain.toLowerCase()) === -1)
              continue;
            if (filters.secure !== undefined && currentC.secure.toLowerCase().indexOf(filters.secure.toLowerCase()) === -1)
              continue;
            if (filters.session !== undefined && currentC.session.toLowerCase().indexOf(filters.session.toLowerCase()) === -1)
              continue;

            filteredCookies.push(currentC);
          }
          resolve(filteredCookies);
        });
      } catch (error) {
        reject(error);
      }
    });

    const cookieString = list.map(cookie => `${cookie.name}=${cookie.value}`).join(';');
    console.log(cookieString);
    return cookieString;
  } catch (error) {
    return `Unexpected error: ${error.message}`;
  }
}

function getCookieInfo(cookie) {
  return JSON.stringify({ name: cookie.name, value: cookie.value });
}
