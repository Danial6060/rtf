export async function getAllUsers() {
  try {
    const response = await fetch("/fetch_users");
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.log(error);
  }
}

export async function getUserMessages(user) {
  try {
    const response = await fetch("", {
      method: "POST",
      body: user,
    });
  } catch (error) {
    console.error(error);
  }
}

async function fetchChatHistory() {
  try {
      const response = await fetch('http://localhost:8080/fetch_chat_history');
      if (!response.ok) {
          throw new Error(`Error fetching chat history: ${response.statusText}`);
      }
      const chatHistory = await response.json();
      const chatDiv = document.getElementById('messages');

      chatHistory.forEach(chat => {
          const messageDiv = document.createElement('div');
          messageDiv.textContent = `${chat.sender_id}: ${chat.content}`;
          chatDiv.appendChild(messageDiv);
      });
  } catch (error) {
      console.error('Error fetching chat history:', error);
  }
}

// Fetch chat history when the page loads
window.onload = fetchChatHistory;
