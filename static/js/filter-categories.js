const dropdown = document.querySelectorAll('.dropdown');

dropdown.forEach(input => input.addEventListener('click', function() {
    input.classList.toggle('active')
}));

function getCurrentURL () {
    return window.location.href
}


const url = getCurrentURL();
const splitUrl = url.split(/[?=]+/);


const container = document.querySelector('.filter-display');

if (splitUrl[1] === undefined) {
    container.style.diplay = "none";
} else {
    const filter = splitUrl[1];
    const parameter = splitUrl[2].replace(/%20/g, " ");
    container.innerText = `${filter} : ${parameter}`;
}

