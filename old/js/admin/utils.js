export function fetchTextHTML(url) {
    return fetch(url).then(res => {
        if (!res.ok) throw new Error("Failed to load " + url);
        return res.text();
    });
}

export function attachSectionToggles() {
    document.querySelectorAll(".collapsible").forEach(section => {
        const title = section.querySelector(".section-title");
        const content = section.querySelector(".section-content");
        title.addEventListener("click", () => {
            content.classList.toggle("collapsed");
        });
    });
}
