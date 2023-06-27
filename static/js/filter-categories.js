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

const container = document.querySelector('.reset');

if (splitUrl[1] === undefined) {
    container.innerHTML = "";
    container.classList.remove("reset");
} else {
    container.innerHTML = "Clear Filter";
    container.classList.add("reset");
}

