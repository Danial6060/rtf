import { registerPage, loginPage, mainPage, messagePage } from "./pages.js";

// render last page before page refresh
document.addEventListener("DOMContentLoaded", () => {
  // get last page rendered from the localStorage
  const currentPage = window.localStorage.getItem("currentPage");

  switch (currentPage) {
    case "register":
      registerPage();
      break;
    case "login":
      loginPage();
      break;
    case "main":
      mainPage();
      break;
    case "chat":
      messagePage();
      break;
    default:
      loginPage();
      break;
  }
});

function connectWebSocket() {
  const ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = function () {
    console.log("Connected to WebSocket");
  };

  ws.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    console.log("Received:", msg);
  };

  ws.onclose = function () {
    console.log("WebSocket connection closed");
  };
}

// document.addEventListener("DOMContentLoaded", function () {
//   connectWebSocket();

//   document
//     .getElementById("register-form")
//     .addEventListener("submit", function (event) {
//       event.preventDefault();
//       registerUser();
//     });

//   document
//     .getElementById("login-form")
//     .addEventListener("submit", function (event) {
//       event.preventDefault();
//       loginUser();
//     });

//   document
//     .getElementById("logout-button")
//     .addEventListener("click", function () {
//       logoutUser();
//     });

//   document
//     .getElementById("post-form")
//     .addEventListener("submit", function (event) {
//       event.preventDefault();
//       createPost();
//     });

//   fetchPosts(); // Ensure this is called to fetch posts when the page loads
// });
