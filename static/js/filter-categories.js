const dropdown = document.querySelectorAll('.dropdown');

dropdown.forEach(input => input.addEventListener('click', function() {
    input.classList.toggle('active')
}));
