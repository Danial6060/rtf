import {
  registerUser,
  loginUser,
  logoutUser,
  checkIfUserLoggedIn,
} from "./auth.js";
import {
  createPost,
  displayComments,
  fetchComments,
  fetchPosts,
} from "./post.js";
import { navigateTo } from "./history.js";
import { getAllUsers } from "./chat.js";

export function registerPage() {
  // Clear all content from the body
  document.body.innerHTML = "";

  document.body.innerHTML = `
        <div class="register-content">
            <form id="register-form">
                <h2>Register</h2>
                <input type="text" name="nickname" placeholder="Nickname" required />
                <input type="number" name="age" placeholder="Age" required />
                <input type="text" name="gender" placeholder="Gender" required />
                <input type="text" name="first_name" placeholder="First Name" required />
                <input type="text" name="last_name" placeholder="Last Name" required />
                <input type="email" name="email" placeholder="Email" required />
                <input type="password" name="password" placeholder="Password" required />
                <button type="submit">Register</button>
            </form>

            <button class="login-button">Go to login</button>
        </div>
    `;

  // render login page
  document.querySelector(".login-button").addEventListener("click", () => {
    navigateTo(loginPage);
  });

  document
    .querySelector("#register-form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      // if no errors were found
      if (await registerUser()) {
        navigateTo(mainPage);
      }
    });

  // setting current page on local storage to make the page persist across page refresh
  window.localStorage.setItem("currentPage", "register");
}

export function loginPage() {
  document.body.innerHTML = "";

  document.body.innerHTML = `
        <div class="login-content">
          <form id="login-form">
          <h2>Login</h2>
          <input
            type="text"
            name="identifier"
            placeholder="Nickname or Email"
            required
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            required
          />
          <button type="submit">Login</button>
        </form>

        <button class="register-button">Go to register</button>
      </div>
  `;

  // render register page
  document.querySelector(".register-button").addEventListener("click", () => {
    navigateTo(registerPage);
  });

  document
    .querySelector("#login-form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();
      if (await loginUser()) {
        navigateTo(mainPage);
      }
    });

  window.localStorage.setItem("currentPage", "login");
}

export async function mainPage() {
  document.body.innerHTML = "";
  document.body.innerHTML = `
        <div id="main-content">
          <div class="nav"></div>
          <form id="post-form">
            <h2>Create Post</h2>
            <input type="text" name="category" placeholder="Category" required />
            <textarea name="content" placeholder="Content" required></textarea>
            <button type="submit">Post</button>
          </form>

          <h2>Posts</h2>

          <div id="posts">
            <!-- Posts will be dynamically loaded here -->
          </div>
      </div>
  `;

  const isLoggedIn = await checkIfUserLoggedIn();

  if (isLoggedIn) {
    document.querySelector(".nav").innerHTML = `
    <button id="logout-button">Logout</button>
    <button id="chat">Chat</div>
    `;

    document.querySelector("#logout-button").addEventListener("click", () => {
      // logout the user
      logoutUser();

      // after logging out render login page (temporary)
      navigateTo(loginPage);
    });

    document.querySelector("#chat").addEventListener("click", () => {
      navigateTo(messagePage);
    });
  } else {
    document.querySelector(
      ".nav"
    ).innerHTML = `<button id="login-button">Login</button>`;

    document.querySelector("#login-button").addEventListener("click", () => {
      navigateTo(loginPage);
    });
  }

  fetchPosts();

  document.querySelector("#post-form").addEventListener("submit", (e) => {
    e.preventDefault();
    createPost();
  });

  window.localStorage.setItem("currentPage", "main");
}

export async function messagePage() {
  const isLoggedIn = await checkIfUserLoggedIn();

  if (!isLoggedIn) {
    navigateTo(loginPage);
    return;
  }

  document.body.innerHTML = `
  <div class="main-content">
   <button id="backButton">Back</button>
    <div class="all-users"></div>
    <div class="selected-user">
      <div id="chatContainer">
        <div id="messages"></div>
        <div id="inputContainer">
          <input type="text" id="messageInput" placeholder="Type your message here">
          <button id="sendButton">Send</button>
        </div>
      </div>
    </div>
  </div>
  `;

  const users = await getAllUsers();
  const userContainer = document.querySelector(".all-users");

  // Populate users
  users.forEach((user) => {
    const div = document.createElement("div");
    div.className = "user";
    div.textContent = user;

    div.addEventListener("click", () => {
      handleUserClick(user);
    });

    userContainer.appendChild(div);
  });

  const sendButton = document.getElementById("sendButton");
  const messageInput = document.getElementById("messageInput");
  const messagesContainer = document.getElementById("messages");
  const backButton = document.getElementById("backButton");

  // WebSocket connection
  const socket = new WebSocket('ws://localhost:8080/ws');

  socket.onopen = () => {
    console.log('WebSocket connection established');
  };

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    displayMessage(msg, 'other-message');
  };

  socket.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  socket.onclose = () => {
    console.log('WebSocket connection closed');
  };

  // Handle sending message
  sendButton.addEventListener("click", () => {
    const message = messageInput.value.trim();
    if (message !== "") {
      sendMessage(message);
      messageInput.value = "";
    }
  });

  // Allow pressing "Enter" to send message
  messageInput.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      sendButton.click();
    }
  });

  // Handle back button click
  document.querySelector("#backButton").addEventListener("click", () => {
    navigateTo(mainPage);
  });

  function sendMessage(content) {
    const msg = {
      SenderID: 1, // Replace with actual sender ID
      ReceiverID: 2, // Replace with actual receiver ID
      Content: content
    };

    socket.send(JSON.stringify(msg));

    displayMessage(msg, 'my-message');
  }

  function displayMessage(msg, className) {
    const messageElement = document.createElement("div");
    messageElement.className = `message ${className}`;
    messageElement.textContent = msg.Content;

    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  window.localStorage.setItem("currentPage", "chat");
}

// making the functions globally accessible to navigate forward and backwards
window.registerPage = registerPage;
window.loginPage = loginPage;
window.mainPage = mainPage;
window.messagePage = messagePage;
