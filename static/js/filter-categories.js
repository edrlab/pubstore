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

const deleteFilter = document.querySelector('.main-content-header span');
const container = document.querySelector('.main-content-header h3');

if (splitUrl[1] === undefined) {
    container.style.diplay = "none";
    container.classList.remove("filter-display");
    deleteFilter.innerHTML = "";
} else {
    container.classList.add("filter-display");
    deleteFilter.innerHTML ='<i class="fa-solid fa-xmark"></i>';
}

