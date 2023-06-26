const dropdown = document.querySelectorAll('.dropdown');

dropdown.forEach(input => input.addEventListener('click', function() {
    input.classList.toggle('active')
}));

function getCurrentURL () {
    return window.location.href
}


const url = getCurrentURL();
const splitUrl = url.split(/[?=]+/);
const decodedUrl = decodeURIComponent(splitUrl[2]);

const container = document.querySelector('.main-content-header h3');

if (splitUrl[1] === undefined) {
    container.style.diplay = "none";
    container.classList.remove("filter-display");
} else {
    container.classList.add("filter-display");
    container.innerHTML = `${decodedUrl} <i class="fa-solid fa-xmark"></i>`;
}

