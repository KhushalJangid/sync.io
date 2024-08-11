var dropzone = document.getElementById('dropzone');
var input = document.getElementById('file');
const uploadButton = document.getElementById('uploadButton');
const uploadModal = document.getElementById('uploadModal');
const closeModal = document.getElementById('closeModal');
const minimizeModal = document.getElementById('minimizeModal');

// Function to show the modal with animation
const showmodal = () => {
    uploadModal.classList.remove('hidden');
    uploadModal.classList.remove('translate-x-full', 'translate-y-full');
}

// Function to hide the modal
closeModal.addEventListener('click', () => {
    uploadModal.classList.add('translate-x-full', 'translate-y-full');
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
    if(input.files.length == 0){
        return;
    }
    event.preventDefault();
    console.log(input.files)
    handleFiles(input.files);
    input.value = null;
    // dropzone.innerText = "Drag and drop files here or click to select files";
    loadTitle();
});
var loadTitle = function (e) {
    if (input.files.length == 1) {
        // dropzone.innerText = input.files.item(0).name;
        dropzone.innerHTML = `
            <p class="mb-2 text-sm text-gray-500 dark:text-white">
            <i class="fa-regular fa-file"></i>&nbsp${input.files.item(0).name}</p>`;
    } else if (input.files.length == 0) {
        // dropzone.innerText = "Drag and drop files here or click to select files";
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
        // dropzone.innerText = `${input.files.length} Files selected`;

        dropzone.innerHTML = `
            <p class="mb-2 text-sm text-gray-500 dark:text-white">
            <i class="fa-regular fa-folder"></i>&nbsp${input.files.item(0).name} and ${input.files.length-1} more files selected</p>`;
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
    // handleFiles(files);
    input.files = files;
    loadTitle(null);
});

dropzone.addEventListener('click', function () {
    document.getElementById('file').click();
});

function handleFiles(files) {
    var progressContainer = document.getElementById('progressContainer');
    progressContainer.innerHTML = ''; // Clear previous progress bars

    // Show the modal
    // $('#progressModal').modal('show');
    showmodal();
    console.log(files)
    Array.from(files).forEach(function (file, index) {
        var formData = new FormData();
        formData.append('file', file);

        // var progressWrapper = document.createElement('div');
        // progressWrapper.className = 'mb-3';

        // var fileInfo = document.createElement('div');
        // fileInfo.className = 'mb-2';
        // fileInfo.textContent = `${file.name} (${(file.size / 1024 / 1024).toFixed(2)} MB)`;

        // var progressBar = document.createElement('div');
        // progressBar.className = 'progress-bar';
        // progressBar.style.width = '0%';
        // progressBar.setAttribute('role', 'progressbar');
        // progressBar.setAttribute('aria-valuenow', '0');
        // progressBar.setAttribute('aria-valuemin', '0');
        // progressBar.setAttribute('aria-valuemax', '100');
        // progressBar.textContent = '0%';

        // var progressBarContainer = document.createElement('div');
        // progressBarContainer.className = 'progress';
        // progressBarContainer.style.height = '15px';
        // progressBarContainer.appendChild(progressBar);

        // var sizeInfo = document.createElement('div');
        // sizeInfo.className = 'text-right';
        // sizeInfo.textContent = `0 / ${(file.size / 1024 / 1024).toFixed(2)} MB`;

        // progressWrapper.appendChild(fileInfo);
        // progressWrapper.appendChild(progressBarContainer);
        // progressWrapper.appendChild(sizeInfo);
        // progressContainer.appendChild(progressWrapper);
        var tile = new Tile(file);
        tile.elements().forEach((e,i,arr)=>{
            progressContainer.append(e);
        })

        var xhr = new XMLHttpRequest();

        xhr.upload.addEventListener('progress', function (event) {
            if (event.lengthComputable) {
                var percentComplete = Math.round((event.loaded / event.total) * 100);
                tile.updateProgress(percentComplete);
                // progressBar.style.width = percentComplete + '%';
                // progressBar.setAttribute('aria-valuenow', percentComplete);
                // progressBar.textContent = percentComplete + '%';
                // sizeInfo.textContent = `${(event.loaded / 1024 / 1024).toFixed(2)} / ${(file.size / 1024 / 1024).toFixed(2)} MB`;
            }
        });

        xhr.open('POST', '/upload', true);
        xhr.send(formData);

        xhr.onload = function () {
            if (xhr.status === 200) {
                // progressBar.classList.add('bg-success');
                tile.markSuccess();
            } else {
                // progressBar.classList.add('bg-danger');
                // progressBar.textContent = 'Upload failed';
                tile.markFailed();
            }
        };
    });
}

class Tile {
    constructor(file) {
        this.tile = document.createElement('div');
        this.tile.className = 'mb-2 flex justify-between items-center';
        this.tile.innerHTML = `
    <div class="flex items-center gap-x-3">
        <!-- File Icon -->
        <span class="size-8 flex justify-center items-center border border-gray-200 dark:border-gray-900 rounded-lg">
                <!--<svg class="shrink-0 size-5 stroke-gray-200 dark:stroke-slate-300"
                    xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                    viewBox="0 0 24 24" fill="none" stroke-width="1.5"
                    stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                    <polyline points="17 8 12 3 7 8"></polyline>
                    <line x1="12" x2="12" y1="3" y2="15"></line>
                </svg>-->
                <i class="fa-regular fa-file"></i>
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
                <!--<svg class="shrink-0 size-4" xmlns="http://www.w3.org/2000/svg" width="24"
                    height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                    stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M3 6h18"></path>
                    <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                    <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
                    <line x1="10" x2="10" y1="11" y2="17"></line>
                    <line x1="14" x2="14" y1="11" y2="17"></line>
                </svg>-->
                <i class="fa-solid fa-xmark"></i>
                <span class="sr-only">Delete</span>
        </button> 
    </div>`;
        this.progress = document.createElement('div');
        this.progress.className = "flex items-center gap-x-3 whitespace-nowrap";
        this.progress.innerHTML = `
    <div class="flex w-full h-2 bg-gray-200 rounded-full overflow-hidden dark:bg-neutral-700"
        role="progressbar" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">
        <div class="flex flex-col justify-center rounded-full overflow-hidden bg-blue-600 text-xs text-white text-center whitespace-nowrap transition duration-500 dark:bg-blue-500"
                style="width: 0%"></div>
    </div>
    <div class="w-6 text-end">
        <span class="text-sm ">0%</span>
    </div>`;
    }
    updateProgress(percent){
        this.progress.getElementsByClassName('transition').item(0).style.width = `${percent}%`;
        this.progress.getElementsByTagName('span').item(0).innerText = `${percent}%`;
    }
    elements(){
        return [this.tile,this.progress];
    }
    cancel(fn){
        this.tile.getElementsByClassName('button').addEventListener('click',fn);
    }
    markSuccess(){
        this.progress.getElementsByClassName('transition').item(0).style.backgroundColor = 'green';
    }
    markFailed(){
        this.progress.getElementsByClassName('transition').item(0).style.backgroundColor = 'red';
    }
}