var dropzone = document.getElementById('dropzone');
var input = document.getElementById('file');
const uploadButton = document.getElementById('uploadButton');
const uploadModal = document.getElementById('uploadModal');
const closeModal = document.getElementById('closeModal');
const minimizeModal = document.getElementById('minimizeModal');

// Function to show the modal with animation
const showmodal = () => {
    uploadModal.classList.remove('hidden');
    uploadModal.classList.remove( 'translate-y-full');
}

// Function to hide the modal
closeModal.addEventListener('click', () => {
    uploadModal.classList.add( 'translate-y-full');
    setTimeout(() => uploadModal.classList.add('hidden'), 300);
});

// Function to minimize the modal
minimizeModal.addEventListener('click', () => {
    if (!uploadModal.classList.contains('minimized')) {
        uploadModal.style.height = '48px';
        uploadModal.classList.add('minimized');
    } else {
        uploadModal.style.height = '';
        uploadModal.classList.remove('minimized');
    }
});

// Add Tailwind CSS animations for sliding in from the bottom right
document.head.insertAdjacentHTML('beforeend', `
<style>
    @keyframes slideIn {
            from {
                transform: translateX(0) translateY(100%);
            }
            to {
                transform: translateX(0) translateY(0);
            }
    }
    .animate-slideIn {
            animation: slideIn 1s ease-out forwards;
    }
</style>
`);

window.onload = function () {
    input.value = null;
}
uploadButton.addEventListener('click', function (event) {
    if (input.files.length == 0) {
        return;
    }
    event.preventDefault();
    console.log(input.files)
    handleFiles(input.files);
    input.value = null;
    loadTitle();
});
var loadTitle = function (e) {
    if (input.files.length == 1) {
        dropzone.innerHTML = `
                <div class="flex flex-col items-center justify-center pt-5 pb-6">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none"
                            class="bi bi-file-earmark w-8 h-8 mb-4 text-gray-500 dark:text-white" viewBox="0 0 16 16">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                        stroke-width="1"
                                d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5z" />
                    </svg>
                    <p class="mb-2 text-sm text-gray-500 dark:text-white"><span
                                class="font-semibold">${input.files.item(0).name}</p>
                </div>`;
    } else if (input.files.length == 0) {
        dropzone.innerHTML = `
        <div class="flex flex-col items-center justify-center pt-5 pb-6">
            <svg class="w-8 h-8 mb-4 text-gray-500 dark:text-white" aria-hidden="true"
                    xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 16">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                        stroke-width="2"
                        d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2" />
            </svg>
            <p class="mb-2 text-sm text-gray-500 dark:text-white"><span
                        class="font-semibold">Click to upload</span> or drag and drop</p>
        </div>`;
    } else {
        dropzone.innerHTML = `
        <div class="flex flex-col items-center justify-center pt-5 pb-6">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none"
                    class="bi bi-folder2 w-8 h-8 mb-4 text-gray-500 dark:text-white" viewBox="0 0 16 16">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                        stroke-width="1"
                        d="M1 3.5A1.5 1.5 0 0 1 2.5 2h2.764c.958 0 1.76.56 2.311 1.184C7.985 3.648 8.48 4 9 4h4.5A1.5 1.5 0 0 1 15 5.5v7a1.5 1.5 0 0 1-1.5 1.5h-11A1.5 1.5 0 0 1 1 12.5zM2.5 3a.5.5 0 0 0-.5.5V6h12v-.5a.5.5 0 0 0-.5-.5H9c-.964 0-1.71-.629-2.174-1.154C6.374 3.334 5.82 3 5.264 3zM14 7H2v5.5a.5.5 0 0 0 .5.5h11a.5.5 0 0 0 .5-.5z" />
            </svg>
            <p class="mb-2 text-sm text-gray-500 dark:text-white"><span
                        class="font-semibold">${input.files.item(0).name} and ${input.files.length - 1} more files selected</p>
        </div>`;
    }
};
input.addEventListener('change', loadTitle);

dropzone.addEventListener('dragover', function (event) {
    event.preventDefault();
    dropzone.classList.add('dragover');
});

dropzone.addEventListener('dragleave', function (event) {
    dropzone.classList.remove('dragover');
});

dropzone.addEventListener('drop', function (event) {
    event.preventDefault();
    dropzone.classList.remove('dragover');
    var files = event.dataTransfer.files;
    input.files = files;
    loadTitle(null);
});

dropzone.addEventListener('click', function () {
    document.getElementById('file').click();
});

function handleFiles(files) {
    var progressContainer = document.getElementById('progressContainer');
    progressContainer.innerHTML = ''; // Clear previous progress bars
    showmodal();
    console.log(files)
    const csrf = document.getElementById('csrf').value;
    Array.from(files).forEach(function (file, index) {
        var formData = new FormData();
        formData.append('file', file);
        var tile = new Tile(file);
        tile.elements().forEach((e, i, arr) => {
            progressContainer.append(e);
        })

        var xhr = new XMLHttpRequest();
        tile.cancel(function (event){
            console.info("Aborting upload");
            xhr.abort();
            tile.markFailed();
        });

        xhr.upload.addEventListener('progress', function (event) {
            if (event.lengthComputable) {
                var percentComplete = Math.round((event.loaded / event.total) * 100);
                tile.updateProgress(percentComplete);
            }
        });

        xhr.open('POST', '/upload', true);
        xhr.setRequestHeader('X-CSRF-TOKEN', csrf);
        xhr.send(formData);

        xhr.onload = function () {
            if (xhr.status === 200) {
                tile.markSuccess();
            } else {
                tile.markFailed();
            }
        };
    });
}

class Tile {
    constructor(file) {
        this.tile = document.createElement('div');
        this.tile.className = 'mt-4 flex justify-between items-center';
        this.tile.innerHTML = `
    <div class="flex items-center gap-x-3">
        <!-- File Icon -->
        <span class="size-8 flex justify-center items-center border border-gray-200 dark:border-gray-900 rounded-lg">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" class="bi bi-file-earmark mr-2" viewBox="0 0 16 16">
                <path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5z"/>
                </svg>
        </span>

        <!-- File Name and Size -->
        <div>
                <p class="text-sm font-medium file-title">${file.name}</p>
                <p class="text-xs dark:text-gray-400 text-neutral-500">${(file.size / 1024 / 1024).toFixed(2)} MB</p>
        </div>
    </div>

    <div class="inline-flex items-center gap-x-2">
        <!-- Pause Button 
        <button type="button"
                class="relative text-gray-200 hover:text-gray-800 focus:outline-none focus:text-gray-800 disabled:opacity-50 disabled:pointer-events-none dark:text-slate-300 dark:hover:text-neutral-200 dark:focus:text-neutral-200">
                <svg class="shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24"
                    height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                    stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect width="4" height="16" x="6" y="4"></rect>
                    <rect width="4" height="16" x="14" y="4"></rect>
                </svg>
                <span class="sr-only">Pause</span>
        </button> -->

        <!-- Delete Button -->
        <button type="button"
                class="relative dark:text-gray-200 focus:outline-none dark:focus:text-gray-800 
                disabled:opacity-50 disabled:pointer-events-none 
                text-slate-300 focus:text-neutral-200">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="red" class="bi bi-x-circle" viewBox="0 0 16 16">
                <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>
                <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/>
                </svg>
                <span class="sr-only">Cancel</span>
        </button> 
    </div>`;
        this.progress = document.createElement('div');
        this.progress.className = "flex items-center gap-x-3 whitespace-nowrap";
        this.progress.innerHTML = `
    <div class="flex h-2 bg-gray-200 rounded-full overflow-hidden dark:bg-neutral-700"
        role="progressbar" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100" style="width: 80%">
        <div class="flex flex-col justify-center rounded-full overflow-hidden bg-blue-600 text-xs text-white text-center whitespace-nowrap transition duration-500 dark:bg-blue-500"
                style="width: 0%"></div>
    </div>
    <div class="w-6 text-end">
        <span class="text-sm ">0%</span>
    </div>`;
    }
    updateProgress(percent) {
        this.progress.getElementsByClassName('transition').item(0).style.width = `${percent}%`;
        this.progress.getElementsByTagName('span').item(0).innerText = `${percent}%`;
    }
    elements() {
        return [this.tile, this.progress];
    }
    cancel(fn) {
        this.tile.getElementsByTagName('button').item(0).addEventListener('click', fn);
    }
    markSuccess() {
        this.progress.getElementsByClassName('transition').item(0).style.backgroundColor = 'green';
    }
    markFailed() {
        this.progress.getElementsByClassName('transition').item(0).style.backgroundColor = 'red';
    }
}