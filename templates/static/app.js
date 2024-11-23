document.addEventListener("DOMContentLoaded", () => {
    const checkboxes = document.querySelectorAll(".subtask-checkbox");

    checkboxes.forEach(checkbox => {
        const listItem = checkbox.closest(".list-group-item");

        // Ensure styles are correctly applied on load
        if (checkbox.checked) {
            listItem.classList.add("bg-success", "text-decoration-line-through");
            listItem.classList.remove("text-dark", "bg-light");
        }

        // Update styles dynamically when checkbox is toggled
        checkbox.addEventListener("change", function () {
            if (this.checked) {
                listItem.classList.add("bg-success", "text-decoration-line-through");
                listItem.classList.remove("text-dark", "bg-light");
            } else {
                listItem.classList.remove("bg-success", "text-decoration-line-through");
                listItem.classList.add("text-dark", "bg-light");
            }
        });
    });
});


document.addEventListener("DOMContentLoaded", () => {
    const doneButtons = document.querySelectorAll(".btn-done");

    doneButtons.forEach(button => {
        button.addEventListener("click", async function (e) {
            e.preventDefault(); // Prevent form submission
            const todoID = this.dataset.todoId;

            try {
                // Send a POST request to mark the todo as done
                const response = await fetch(`/todos/${todoID}/markdone`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (response.ok) {
                    // Find the parent card
                    const todoCard = this.closest(".col-md-6");
                    const doneSection = document.querySelector("#done-tickets");

                    // Move the todo card to the "Done Tickets" section
                    doneSection.appendChild(todoCard);

                    // Update the card's styling to indicate completion
                    todoCard.querySelector(".badge").classList.remove("bg-danger", "bg-warning");
                    todoCard.querySelector(".badge").classList.add("bg-success");
                    todoCard.querySelector(".badge").textContent = "Completed";

                    // Remove the "Done?" button
                    this.parentElement.removeChild(this);
                } else {
                    console.error("Failed to mark todo as done");
                }
            } catch (error) {
                console.error("Error:", error);
            }
        });
    });
});
