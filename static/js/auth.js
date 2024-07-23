export async function registerUser() {
  const data = new FormData(document.getElementById("register-form"));
  try {
    const response = await fetch("/register", {
      method: "POST",
      body: data,
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    return true;
  } catch (error) {
    console.error(error);
    return false;
  }
}

export async function loginUser() {
  const data = new FormData(document.querySelector("#login-form"));

  try {
    const response = await fetch("/login", {
      method: "POST",
      body: data,
      credentials: "include", // Ensure cookies are included in the request
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    return true;
  } catch (error) {
    console.error(error);
    return false;
  }
}

export async function logoutUser() {
  try {
    const response = await fetch("/logout", {
      method: "POST",
    });

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
  } catch (error) {
    console.error(error);
  }
}

export async function checkIfUserLoggedIn() {
  try {
    const response = await fetch("/is_logged_in", {
      method: "GET",
      credentials: "include", // Ensure cookies are included in the request
    });

    if (response.status === 401) {
      // throw an error with the error message from the backend
      return false;
    }

    if (!response.ok) {
      throw new Error(await response.text());
    }

    return true;
  } catch (error) {
    console.error(error);
    return false;
  }
}
