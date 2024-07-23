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
