const modalWindow = document.getElementById('modal-window');
const modal = document.getElementById('modal');
const modalBackGround = document.querySelector('.modal-backdrop');

let modalState = null;

const purchaseButton = document.getElementById("buy");
const loanButton = document.getElementById("loan");

const buyForm = document.getElementById("buyForm");
const loanForm = document.getElementById("loanForm");

const submitButtonBuy = document.getElementById("submitBuy");
const submitButtonLoan = document.getElementById("submitLoan");

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

console.log({
  purchaseButton,
  loanButton,
  submitButtonBuy,
  submitButtonLoan,
  modalWindow,
  modal,
  buyForm,
  loanForm
});

loanButton.addEventListener('click', (e) => {
    e.preventDefault();
    createModal()
    buyForm.style.display = "none";
    loanForm.style.display = "flex";
})


submitButtonBuy.addEventListener('click', (e) => {
    setTimeout(() => {
        window.location.reload();
    }, 1000)
});

submitButtonLoan.addEventListener('click', (e) => {
    setTimeout(() => {
        window.location.reload();
    }, 1000)
});

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
