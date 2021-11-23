let apiUrl="/api/thumbnails/"

window.onscroll = fillGallery;
window.onkeydown = actUponKeyPressed;
const RootGalleryElement = document.getElementById("gallery");
const RootOriginalViewElement = document.getElementById("original_view");
var imageCount = 0;

async function main() {
	getNextImages()
}

function actUponKeyPressed(e) {
	e = e || window.event;

	switch(e.keyCode) {
		case 27:
			// escape
			hideOriginal();
			break;
		case 38: 
			// arrow up
			hideOriginal();
			break;
        case 37:
        	console.log("left arrow");
        	break;
        case 39:
        	console.log("right arrow");        
            break;
    }

}

async function fillGallery() {
	const scrolled = document.body.scrollTop ? document.body.scrollTop : document.documentElement.scrollTop;
	const windowHeight = window.innerHeight;
	const overflowHeight = scrolled + windowHeight + 400;
	const lastImagePosition = Math.floor(imageCount / 4)*200;
	if (lastImagePosition < overflowHeight) { getNextImages() }
}

async function getNextImages() {
	await getImages(imageCount, imageCount+40);
	imageCount = imageCount + 40;
}

async function getImages(from,to) {
	const response = await fetch(apiUrl+from+"/"+to)
	const result = await response.json()
	for (let imageIndex in result.thumbnails) {
		appendImage(result.thumbnails[imageIndex], RootGalleryElement);
	}
	return result
}

function appendImage(imageUrl, rootGalleryElement) {
    const elem = document.createElement("div");
    elem.classList.add("picture");
    elem.style.backgroundImage = "url(/thumbnails/"+imageUrl+")";
    elem.addEventListener("click", showOriginal);
    elem.setAttribute("picture_name", imageUrl);
	rootGalleryElement.appendChild(elem);
}

async function showOriginal(clickEvent) {
	const target = clickEvent.target;
	const originalUrl = target.attributes["picture_name"].value;
	console.log(originalUrl);
	const original = document.createElement("div");
	original.style.width = "100%";
	original.style.height = "100%";
	original.style.position = "absolute";
	original.style.display = "block";
	original.style.backgroundImage = "url(/originals/"+originalUrl+")";
	original.style.backgroundPosition = "center center";
	original.style.backgroundRepeat = "no-repeat";
	original.style.backgroundSize = "contain";
	original.style.backgroundColor = "#222";
	original.addEventListener("click", hideOriginal);

	hideOriginal()
	RootOriginalViewElement.appendChild(original);

	const earlierElement = target.previousElementSibling;
	const laterElement = target.nextElementSibling;
	if (earlierElement) {
		const showEarlier = document.createElement("div");
		showEarlier.style.position = "fixed";
		showEarlier.style.width = "35%";
		showEarlier.style.height = "100%";
		showEarlier.style.left = 0;
		showEarlier.style.top = 0;
		showEarlier.addEventListener("mouseover", (e) => showEarlier.classList.add("hoverShowEarlier"));
		showEarlier.addEventListener("mouseout", (e) => showEarlier.classList.remove("hoverShowEarlier"));
		showEarlier.addEventListener("click", (e) => showOriginal({ target: earlierElement }));
		RootOriginalViewElement.appendChild(showEarlier);
	}
	if (laterElement) {
		const showLater = document.createElement("div");
		showLater.style.position = "fixed";
		showLater.style.width = "35%";
		showLater.style.height = "100%";
		showLater.style.left = "65%";
		showLater.style.top = 0;
		showLater.style.borderLeftStyle = "";
		showLater.addEventListener("mouseover", (e) => showLater.classList.add("hoverShowLater"));
		showLater.addEventListener("mouseout", (e) => showLater.classList.remove("hoverShowLater"));
		showLater.addEventListener("click", (e) => showOriginal({ target: laterElement }));
		RootOriginalViewElement.appendChild(showLater);
	}
	RootOriginalViewElement.style.display = "block";
}

async function hideOriginal() {
	RootOriginalViewElement.style.display = "none";
	while (RootOriginalViewElement.firstChild) {
    	RootOriginalViewElement.removeChild(RootOriginalViewElement.lastChild);
  	}
}

function photograf_main() {
	main().catch((error) => console.error(error))
}
