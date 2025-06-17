let currentPage = 1;
let editingId = null;
let hasNextPage = true;
let hasPrevPage = false;


function formatDate(isoDateString) {
    const date = new Date(isoDateString);
    return date.toLocaleString("ru-RU", {
        timeZone: "UTC",
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

function logout() {
    fetch("/api/logout", {
        method: "POST",
    }).then(res => {
        if (res.redirected) {
            window.location.href = res.url;
        }
    });
}

async function updateListings() {
    const filter = document.getElementById('filter').value;
    const res = await fetch(`/api/listings?page=${currentPage}&filter=${filter}`, {
    });
    const data = await res.json();
    const tbody = document.getElementById('listings');
    tbody.innerHTML = '';

    if (data === null) {
        hasNextPage = false;
        if (currentPage > 1) {
            currentPage--;
             await updateListings();
            return;
        }
        const tr = document.createElement('tr');
        tr.innerHTML = `<td colspan="9" style="text-align:center; padding: 10px;">Нет объявлений</td>`;
        tbody.appendChild(tr);
    } else {
        hasNextPage = data.length === 10;
        hasPrevPage = currentPage > 1;

        data.forEach((item, index) => {

            const tr = document.createElement('tr');
            const itemNumber = (currentPage - 1) * 10 + index + 1; // глобальный номер
            tr.innerHTML = `
       <td>${itemNumber}</td>
      <td><strong>${item.Name}</strong></td>
      <td>${item.Typel}</td>
      <td>${item.Description}</td>
      <td>${item.Status}</td>
      <td>${item.Price.toLocaleString("ru-RU", {
                style: "currency",
                currency: "RUB"
            })}</td>
      <td>${item.City}</td>
      <td>${formatDate(item.Date_created)}</td>
      
      <td><button class="action-button edit-btn" onclick=openModal(${JSON.stringify(item)})>Изменить</button> <button class="action-button delete-btn" onclick="deleteListing(${item.ID})">Удалить</button></td>
    `;
            tbody.appendChild(tr);
        });
    }
    document.getElementById('page-info').textContent = currentPage;
    document.getElementById("prev-button").disabled = !hasPrevPage;
    document.getElementById("next-button").disabled = !hasNextPage;

}

async function deleteListing(id) {
    if (!confirm("Вы уверены, что хотите удалить объявление?")) return;
    await fetch(`/api/listings/${id}`, {
        method: 'DELETE',
    });
    await updateListings();
    await updateAnalytics();
    showToast("Объявление удалено", "#f87171");
}

function openModal(listing = null) {
    // document.getElementById('modal').classList.remove('hidden');
    const saveBtn = document.getElementById('modal-save-button');
    const titleInput = document.getElementById('modal-title');
    const priceInput = document.getElementById('modal-price');
    const cityInput = document.getElementById('modal-city');
    const typeInput = document.getElementById('modal-type');
    const descriptionInput = document.getElementById('modal-description');
    const statusInput = document.getElementById('modal-status');
    const label = document.getElementById('modal-title-label')

    if (listing) {
        // режим редактирования
        editingId = listing.ID;
        typeInput.value = listing.Typel
        descriptionInput.value = listing.Description;
        statusInput.value = listing.Status;
        titleInput.value = listing.Name;
        priceInput.value = listing.Price;
        cityInput.value = listing.City;
        saveBtn.textContent = "Сохранить";
        label.textContent = "Редактировать объявление";
    } else {
        // режим добавления
        editingId = null;
        titleInput.value = "";
        priceInput.value = "";
        cityInput.value = "";
        saveBtn.textContent = "Создать";
        label.textContent = "Новое объявление";
    }
    saveBtn.onclick = () => {
        if (editingId) {
            updateListing(editingId);
        } else {
            createListing();
        }
    }
    document.getElementById('modal').classList.remove('hidden');
}

async function updateAnalytics() {
    const res = await fetch('/api/analytics', {});
    const data = await res.json();

    const topCitiesText = data.top_cities
        .map(c => `${c.city} (${c.count})`)
        .join(", ");

    const container = document.getElementById("analytics");
    container.innerHTML = `
    <strong>Объявлений:</strong> ${data.total_listings} &nbsp;|&nbsp;
    <strong>Средняя цена:</strong> ${data.avg_price.toLocaleString()} ₽ &nbsp;|&nbsp;
    <strong>Топ города:</strong> ${topCitiesText}
  `;
}

function nextPage() {
    // currentPage++;
    // updateListings();
    if (!hasNextPage) return;
    currentPage++;
    updateListings();
}

function prevPage() {
    // if (currentPage > 1) {
    //     currentPage--;
    //     updateListings();
    // }
    if (currentPage <= 1) return;
    currentPage--;
    updateListings();
}

function closeModal() {
    document.getElementById('modal').classList.add('hidden');
    document.getElementById('modal-title').value = '';
    document.getElementById('modal-price').value = '';
    document.getElementById('modal-city').value = '';
    document.getElementById('modal-type').value = '';
    document.getElementById('modal-description').value = '';
    document.getElementById('modal-status').value = '';
}

async function createListing() {
    const title = document.getElementById('modal-title').value;
    const type = document.getElementById('modal-type').value;
    const description = document.getElementById('modal-description').value;
    const status = document.getElementById('modal-status').value;
    const price = parseInt(document.getElementById('modal-price').value);
    const city = document.getElementById('modal-city').value;

    await fetch('/api/listings', {
        method: 'POST',
        // headers: {
        //     'Content-Type': 'application/json',
        //     'Authorization': `Bearer ${getToken()}`
        // },
        body: JSON.stringify({ title, type, description, status, price, city })
    }).then(res => {
        if (res.ok){
            showToast("успешно добавленно", "#22c55e")
        }
        }
    );

    closeModal();
    currentPage = 1;
    updateListings();
    updateAnalytics();
}

async function loadCities() {
    const res = await fetch("/api/cities", {
        // headers: { Authorization: `Bearer ${getToken()}` }
    });
    const cities = await res.json();
    if (cities===null) return;
    const select = document.getElementById("filter");
    cities.forEach(city => {
        const option = document.createElement("option");
        option.value = city;
        option.textContent = city;
        select.appendChild(option);
    });
}

async function updateListing(id) {
    const title = document.getElementById('modal-title').value;
    const type = document.getElementById('modal-type').value;
    const description = document.getElementById('modal-description').value;
    const status = document.getElementById('modal-status').value;
    const price = parseInt(document.getElementById('modal-price').value);
    const city = document.getElementById('modal-city').value;

    await fetch(`/api/listings/${id}`, {
        method: 'PUT',
        // headers: {
        //     'Content-Type': 'application/json',
        //     'Authorization': `Bearer ${getToken()}`
        // },
        body: JSON.stringify({ title, type, description, status, price, city })
    }).then(res => {
            if (res.ok){
                showToast("успешно Изменено", "#22c55e")
            }
        }
    );

    closeModal();
    updateListings();
    updateAnalytics();
}

function showToast(text, background = "#2dd4bf", duration = 2000) {
    const toast = document.getElementById('toast');
    toast.textContent = text;
    toast.style.backgroundColor = background;

    toast.classList.remove('hidden');
    toast.classList.add('show');

    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.classList.add('hidden'), 300);
    }, duration);
}



let allUsers = [];
let allListings = [];

async function fetchAdminData() {
    const users = await fetch('/api/admin/users').then(res => res.json());
    const listings = await fetch('/api/admin/listings').then(res => res.json());

    allUsers = users;
    allListings = listings;

    renderUsers();
    renderUserFilter();
    renderListings();
}

function renderUsers() {
    const usersEl = document.getElementById("users");
    usersEl.innerHTML = "";
    allUsers.forEach(u => {
        usersEl.innerHTML += `
          <tr>
            <td>${u.id}</td>
            <td>${u.name}</td>
            <td>${u.login}</td>
            <td>${u.total}</td>
            <td>${u.role}</td>
            <td>
              <button onclick="setRole(${u.id}, '${u.role === 'admin' ? 'agent' : 'admin'}')">
                Сделать ${u.role === 'admin' ? 'агентом' : 'админом'}
              </button>
              <button onclick="deleteUser(${u.id})" class="action-button delete-btn">Удалить</button>
            </td>
          </tr>
        `;
    });
}

async function setRole(userId, role) {
    await fetch('/api/admin/set-role', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, role })
    });
    showToast("Роль обновлена");
    fetchAdminData();
}

async function deleteUser(userId) {
    if (!confirm("Удалить пользователя?")) return;
    await fetch('/api/admin/delete-user', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId })
    });
    showToast("Пользователь удалён", "#f87171");
    fetchAdminData();
}

function filterUsers() {
    const term = document.getElementById("user-search").value.toLowerCase();
    const rows = document.querySelectorAll("#users tr");
    rows.forEach(row => {
        const cells = row.querySelectorAll("td");
        if (cells.length === 0) return;
        const name = cells[1].textContent.toLowerCase();
        const username = cells[2].textContent.toLowerCase();

        const matches = name.includes(term) || username.includes(term);
        row.style.display = matches ? "" : "none";
    });
}

function renderUserFilter() {
    const select = document.getElementById("user-filter");
    select.innerHTML = '<option value="">Все агенты</option>';
    allUsers.forEach(u => {
        const option = document.createElement("option");
        option.value = u.name;
        option.textContent = u.name;
        select.appendChild(option);
    });
}

function renderListings() {
    const selected = document.getElementById("user-filter").value;

    const listingsEl = document.getElementById("listings");
    listingsEl.innerHTML = "";

    allListings
        .filter(l => !selected || l.Agent === selected)
        .forEach(l => {
            listingsEl.innerHTML += `
            <tr>
              <td>${l.ID}</td>
              <td>${l.Name}</td>
              <td>${l.Typel}</td>
              <td>${l.Description}</td>
              <td>${l.Status}</td>
              <td>${l.Price.toLocaleString("ru-RU", {
                style: "currency",
                currency: "RUB"
            })}</td>
                <td>${l.City}</td>
                <td>${new Date(l.Date_created).toLocaleDateString()}</td>
                <td>${l.Agent}</td>
                <td>
                  <button onclick="deleteAdminListing(${l.ID})" class="delete-btn action-button" ">Удалить</button>
                </td>
            </tr>
          `;
        });
}

async function deleteAdminListing(id) {
    if (!confirm("Вы уверены, что хотите удалить объявление?")) return;
    await fetch(`/api/listings/${id}`, {
        method: 'DELETE',
    });
    showToast("Объявление удалено", "#f87171");
    fetchAdminData();
}