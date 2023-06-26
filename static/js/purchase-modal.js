const modalWindow = document.querySelector('#modal-window');
const modal = document.querySelector('.modal');
const modalBackGround = document.querySelector('.modal-backdrop');
let modalState = null;


const purchaseButton = document.querySelector("#buy");
const loanButton = document.querySelector("#loan");

const buyForm = document.getElementById("buyForm");
const loanForm = document.getElementById("loanForm");

const submitButton = document.querySelector(".modal-loan-buttons button[type='submit']");
const backToPresentationButton = document.querySelector(".back-button");

let form = document.querySelector(".modal-form-options");


const createModal = () => {
    modalWindow.style.display = 'block';
    modalWindow.style.opacity = "1";
    modalBackGround.style.opacity = "0.3";
    modalBackGround.style.zIndex = "1050";
    modalState = modal;
}

purchaseButton.addEventListener('click', (e) => {
    e.preventDefault();
    createModal();
    buyForm.style.display = "flex";
    loanForm.style.display = "none";
})

loanButton.addEventListener('click', (e) => {
    e.preventDefault();
    createModal()
    buyForm.style.display = "none";
    loanForm.style.display = "flex";
})

submitButton.addEventListener('click', (e) => {
    closeModal(e)
});

backToPresentationButton.addEventListener("click", (e) => {
    e.preventDefault();
    closeModal(e)
})


const closeModal = (e) => {
    modalWindow.style.display = 'none';
    modalWindow.style.opacity = "0";
    modalBackGround.style.opacity = "0";
    modalBackGround.style.zIndex = "-1";
    modalState = null;
    buyForm.style.display = "none";
    loanForm.style.display = "none";
}


window.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' || e.key === 'Esc') {
        closeModal(e)
    }
});

modalWindow.addEventListener('click', function(e) {
    closeModal(e);
});

modal.addEventListener('click', (e) => {
    e.stopPropagation();
})