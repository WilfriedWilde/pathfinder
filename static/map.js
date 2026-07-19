const svgObject = document.getElementById("map");

svgObject.addEventListener("load", () => {
    const svg = svgObject.contentDocument.querySelector("svg");
    const svgWindow = svg.ownerDocument.defaultView;
    const viewBox = svg.viewBox.baseVal;

    const prevButton = document.getElementById("prev");
    const nextButton = document.getElementById("next");
    const turnLabel = document.getElementById("turnLabel");

    const TRAIN_X = 10;
    const TRAIN_Y = 10;
    const LABEL_X = 22;
    const LABEL_Y = 3;

    let dragging = false;
    let lastX = 0;
    let lastY = 0;
    let currentTurn = 0;

    svg.addEventListener("wheel", e => {
        e.preventDefault();

        const rect = svg.getBoundingClientRect();

        const mouseX = (e.clientX - rect.left) / rect.width;
        const mouseY = (e.clientY - rect.top) / rect.height;

        const oldWidth = viewBox.width;
        const oldHeight = viewBox.height;

        if (e.deltaY < 0) {
            viewBox.width /= 1.2;
            viewBox.height /= 1.2;
        } else {
            viewBox.width *= 1.2;
            viewBox.height *= 1.2;
        }

        viewBox.x += (oldWidth - viewBox.width) * mouseX;
        viewBox.y += (oldHeight - viewBox.height) * mouseY;
    }, { passive: false });

    svg.addEventListener("mousedown", e => {
        e.preventDefault();
        dragging = true;
        lastX = e.clientX;
        lastY = e.clientY;
    });

    svgWindow.addEventListener("mouseup", () => {
        dragging = false;
    });

    svgWindow.addEventListener("mousemove", e => {
        if (!dragging) return;

        const dx = e.clientX - lastX;
        const dy = e.clientY - lastY;
        const svgRect = svg.getBoundingClientRect();

        viewBox.x -= dx * (viewBox.width / svgRect.width);
        viewBox.y -= dy * (viewBox.height / svgRect.height);

        lastX = e.clientX;
        lastY = e.clientY;
    });

    function updateButtons() {
        turnLabel.textContent = `Turn ${currentTurn + 1} / ${turns.length}`;

        prevButton.disabled = currentTurn === 0;
        nextButton.disabled = currentTurn === turns.length - 1;
    }

    function showTurn(turnIndex) {
        svg.querySelectorAll(".train").forEach(train => {
            train.setAttribute("visibility", "hidden");
        });

        const entries = turns[turnIndex].split(" ");

        for (const entry of entries) {
            const [trainName, stationName] = entry.split("-");

            const station = svg.getElementById(stationName);
            const train = svg.getElementById(trainName);

            const x = station.cx.baseVal.value;
            const y = station.cy.baseVal.value;

            train.setAttribute("visibility", "visible");

            train.querySelector("circle").setAttribute("cx", x + TRAIN_X);
            train.querySelector("circle").setAttribute("cy", y + TRAIN_Y);

            train.querySelector("text").setAttribute("x", x + LABEL_X);
            train.querySelector("text").setAttribute("y", y + LABEL_Y);
        }
    }

    prevButton.onclick = () => {
        if (currentTurn === 0) return;

        currentTurn--;
        showTurn(currentTurn);
        updateButtons();
    };

    nextButton.onclick = () => {
        if (currentTurn === turns.length - 1) return;

        currentTurn++;
        showTurn(currentTurn);
        updateButtons();
    };

    window.addEventListener("keydown", e => {
        if (e.key === "ArrowLeft") {
            e.preventDefault();
            prevButton.click();
        } else if (e.key === "ArrowRight") {
            e.preventDefault();
            nextButton.click();
        }
    });

    showTurn(currentTurn);
    updateButtons();
});
