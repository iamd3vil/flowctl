// --- Global Toast Helper ---
function showToast(message, type = "info") {
  const container = document.getElementById("toast-container");
  if (!container) return;
  const toastId = `toast-${Date.now()}`;

  const alertTypes = {
    info: "alert-info",
    success: "alert-success",
    error: "alert-error",
    warning: "alert-warning",
  };
  const alertClass = alertTypes[type] || "alert-info";

  const toastHtml = `
        <div id="${toastId}" class="alert ${alertClass} shadow-lg animate-fade-in-down">
            <span>${message}</span>
        </div>`;
  container.insertAdjacentHTML("beforeend", toastHtml);

  setTimeout(() => {
    document.getElementById(toastId)?.remove();
  }, 5000);
}

// --- API Helper ---
const api = {
  async request(method, url, data = null) {
    const options = {
      method,
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        // 'Authorization': `Bearer ${localStorage.getItem('authToken')}` // Example for token-based auth
      },
    };
    if (data) {
      options.body = JSON.stringify(data);
    }

    try {
      const response = await fetch(url, options);
      if (!response.ok) {
        const errorData = await response
          .json()
          .catch(() => ({ message: response.statusText }));
        throw new Error(errorData.message || "An unknown error occurred.");
      }
      if (response.status === 204) {
        // No Content
        return null;
      }
      return response.json();
    } catch (error) {
      console.error(`API ${method} Error:`, error);
      showToast(error.message, "error");
      throw error;
    }
  },
  get(url) {
    return this.request("GET", url);
  },
  post(url, data) {
    return this.request("POST", url, data);
  },
  put(url, data) {
    return this.request("PUT", url, data);
  },
  delete(url) {
    return this.request("DELETE", url);
  },
};

// --- Alpine.js Components ---
document.addEventListener("alpine:init", () => {
  Alpine.data("appRouter", () => ({
    currentPage: "",
    init() {
      this.loadPage(window.location.pathname);
      window.onpopstate = () => this.loadPage(window.location.pathname);
    },
    async loadPage(path) {
      const contentEl = document.getElementById("main-content");
      if (!contentEl) return;

      const routes = {
        "/": "/pages/flows.html",
        "/login": "/pages/login.html",
        "/users": "/pages/users.html",
        "/groups": "/pages/groups.html",
        // Add other static routes here
      };
      let pageUrl = routes[path];

      if (!pageUrl) {
        const flowMatch = path.match(/^\/flow\/([a-zA-Z0-9_-]+)$/);
        if (flowMatch) {
          pageUrl = "/pages/flow_input.html";
          window.currentEntityId = flowMatch[1];
        } else {
          pageUrl = "/pages/404.html";
        }
      }

      try {
        const response = await fetch(pageUrl);
        if (!response.ok) throw new Error(`Page not found at ${pageUrl}`);
        contentEl.innerHTML = await response.text();
        this.currentPage = path;
        window.scrollTo(0, 0);
      } catch (e) {
        console.error(e);
        contentEl.innerHTML = `<div class="text-center p-16"><h1 class="text-4xl font-bold">Error</h1><p>${e.message}</p></div>`;
      }
    },
    navigate(path) {
      if (path === this.currentPage) return;
      history.pushState(null, "", path);
      this.loadPage(path);
    },
    logout() {
      api
        .post("/logout")
        .then(() => {
          showToast("Signed out successfully.", "success");
          this.navigate("/login");
        })
        .catch(() => this.navigate("/login")); // Navigate even if logout fails
    },
  }));

  Alpine.data("loginPage", () => ({
    credentials: { username: "", password: "" },
    error: "",
    signIn() {
      this.error = "";
      api
        .post("/login", this.credentials)
        .then(() => {
          showToast("Login successful!", "success");
          document
            .querySelector('[x-data^="appRouter"]')
            .__x.$data.navigate("/");
        })
        .catch((err) => {
          this.error = err.message || "Invalid username or password.";
        });
    },
  }));

  Alpine.data("flowsPage", () => ({
    flows: [],
    init() {
      api
        .get("/api/flows")
        .then((data) => (this.flows = data || []))
        .catch((err) => console.error(err));
    },
    navigate(path) {
      document.querySelector('[x-data^="appRouter"]').__x.$data.navigate(path);
    },
  }));

  Alpine.data("flowInputPage", () => ({
    flow: { meta: {}, inputs: [] },
    formValues: {},
    validationErrors: {},
    init() {
      const flowId = window.currentEntityId;
      if (!flowId) return;
      api.get(`/api/flows/${flowId}`).then((data) => {
        this.flow = data;
        this.flow.inputs.forEach((input) => (this.formValues[input.name] = ""));
      });
    },
    triggerFlow() {
      this.validationErrors = {};
      const flowId = window.currentEntityId;
      api
        .post(`/api/trigger/${flowId}`, this.formValues)
        .then((response) => {
          showToast("Flow triggered successfully!", "success");
          // Optional: navigate to a results page
          // document.querySelector('[x-data^="appRouter"]').__x.$data.navigate(`/results/${response.executionId}`);
        })
        .catch((err) => {
          if (err.errors) this.validationErrors = err.errors;
        });
    },
  }));

  Alpine.data("userManagementPage", () => ({
    users: [],
    allGroups: [],
    searchTerm: "",
    isEditMode: false,
    currentUser: {},
    init() {
      this.loadUsers();
      this.loadGroups();
    },
    loadUsers() {
      api
        .get(`/api/users?search=${this.searchTerm}`)
        .then((data) => (this.users = data || []));
    },
    loadGroups() {
      api.get("/api/groups").then((data) => (this.allGroups = data || []));
    },
    openAddModal() {
      this.isEditMode = false;
      this.currentUser = { name: "", username: "" };
      document.getElementById("user_modal").showModal();
    },
    openEditModal(user) {
      this.isEditMode = true;
      this.currentUser = JSON.parse(JSON.stringify(user));
      this.currentUser.groupIds = user.groups.map((g) => g.id);
      document.getElementById("user_modal").showModal();
    },
    closeModal() {
      document.getElementById("user_modal").close();
    },
    saveUser() {
      const promise = this.isEditMode
        ? api.put(`/api/users/${this.currentUser.id}`, this.currentUser)
        : api.post("/api/users", this.currentUser);
      promise
        .then(() => {
          showToast(
            `User ${this.isEditMode ? "updated" : "created"}.`,
            "success",
          );
          this.loadUsers();
          this.closeModal();
        })
        .catch((err) => {});
    },
    deleteUser(userId) {
      if (confirm("Are you sure?")) {
        api
          .delete(`/api/users/${userId}`)
          .then(() => {
            showToast("User deleted.", "success");
            this.loadUsers();
          })
          .catch((err) => {});
      }
    },
  }));

  Alpine.data("groupManagementPage", () => ({
    groups: [],
    searchTerm: "",
    newGroup: { name: "", description: "" },
    init() {
      this.loadGroups();
    },
    loadGroups() {
      api
        .get(`/api/groups?search=${this.searchTerm}`)
        .then((data) => (this.groups = data || []));
    },
    openAddModal() {
      this.newGroup = { name: "", description: "" };
      document.getElementById("group_modal").showModal();
    },
    closeModal() {
      document.getElementById("group_modal").close();
    },
    saveGroup() {
      api
        .post("/api/groups", this.newGroup)
        .then(() => {
          showToast("Group created successfully.", "success");
          this.loadGroups();
          this.closeModal();
        })
        .catch((err) => {});
    },
    deleteGroup(groupId) {
      if (confirm("Are you sure?")) {
        api
          .delete(`/api/groups/${groupId}`)
          .then(() => {
            showToast("Group deleted.", "success");
            this.loadGroups();
          })
          .catch((err) => {});
      }
    },
  }));
});
