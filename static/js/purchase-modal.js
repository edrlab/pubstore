const modalWindow = document.querySelector('#modal-window');
const modal = document.querySelector('.modal');
const modalBackGround = document.querySelector('.modal-backdrop');
let modalState = null;
let cover = document.querySelector('.pub-left-side img');
let title = document.querySelector('.pub-right-side h2');


const purchaseButton = document.querySelector("#buy");
const loanButton = document.querySelector("#loan");
const buyForm = document.getElementById("buyForm");
const loanForm = document.getElementById("loanForm");
const SubmitButton = document.querySelector(".modal-loan-buttons button[type='submit']");
const backToPresentationButton = document.querySelector(".back-button");
const loanOptions = document.querySelector(".select-loan-dates");
const errorMessage = document.querySelector(".modal-form-options span");

let form = document.querySelector(".modal-form-options");

const date = new Date();
let currentDay= String(date.getDate()).padStart(2, '0');
let currentMonth = String(date.getMonth()+1).padStart(2,"0");
let currentYear = date.getFullYear();
let currentDate = `${currentYear}-${currentMonth}-${currentDay}`;
let testDate = new Date().toISOString();
console.log(testDate);

let myBooksArray = [];

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

    backToPresentationButton.addEventListener("click", (e) => {
        e.preventDefault();
        closeModal(e)
    })

    loanButton.addEventListener('click', (e) => {
        e.preventDefault();
        createModal()
        buyForm.style.display = "none";
        loanForm.style.display = "flex";
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