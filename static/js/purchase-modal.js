const modalWindow = document.querySelector('#modal-window');
const modal = document.querySelector('.modal');
const modalBackGround = document.querySelector('.modal-backdrop');
let modalState = null;


const purchaseButton = document.querySelector("#buy");
const loanButton = document.querySelector("#loan");

const buyForm = document.getElementById("buyForm");
const loanForm = document.getElementById("loanForm");

const submitButtonBuy = document.querySelector(".modal-loan-buttons button[value='Buy']");
const submitButtonLoan = document.querySelector(".modal-loan-buttons button[value='Loan']");
const backToPresentationButtonBuy = document.querySelector("#backButtonBuy");
const backToPresentationButtonLoan = document.querySelector("#backButtonLoan");

let form = document.querySelector(".modal-form-options");

let startDateLocal = document.querySelector('#startDateLocal');
let endDateLocal = document.querySelector('#endDateLocal');
let startDateISO = document.querySelector('#startDate');
let endDateISO = document.querySelector('#endDate');


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


submitButtonBuy.addEventListener('click', (e) => {
    closeModal(e)
});

submitButtonLoan.addEventListener('click', (e) => {
    closeModal(e)
});

backToPresentationButtonBuy.addEventListener("click", (e) => {
    e.preventDefault();
    closeModal(e);
})

backToPresentationButtonLoan.addEventListener("click", (e) => {
    e.preventDefault();
    closeModal(e);
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
});


startDateLocal.addEventListener('change', (e) => {
    startDateISO.value = new Date(startDateLocal.value).toISOString();
});

endDateLocal.addEventListener('change', (e) => {
    endDateISO.value = new Date(endDateLocal.value).toISOString();
});