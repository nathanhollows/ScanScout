function formatDateTimeBasedOnProximity(utcDate) {
	const now = new Date();
	const localDate = new Date(utcDate);

	const oneDayMs = 24 * 60 * 60 * 1000; // Milliseconds in one day
	const oneWeekMs = 7 * oneDayMs;

	const diff = now - localDate;

	const isToday = localDate.toDateString() === now.toDateString();
	const isThisWeek = diff < oneWeekMs && !isToday;

	let formattedDate;

	if (isToday) {
		// Show only the time if the date is today
		formattedDate = new Intl.DateTimeFormat(navigator.language, {
			hour: 'numeric',
			minute: 'numeric'
		}).format(localDate);
	} else if (isThisWeek) {
		// Show weekday name and time if the date is within the last week
		formattedDate = new Intl.DateTimeFormat(navigator.language, {
			weekday: 'long',
			hour: 'numeric',
			minute: 'numeric'
		}).format(localDate);
	} else {
		// Show full date if the date is older than a week
		formattedDate = new Intl.DateTimeFormat(navigator.language, {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: 'numeric'
		}).format(localDate);
	}

	return formattedDate;
}

// Function to convert UTC to local time and update the element's content
function convertUTCtoLocal() {
	const elements = document.querySelectorAll('.convert-time');

	elements.forEach((el) => {
		const utcDateStr = el.getAttribute('data-datetime');
		const localDateFormatted = formatDateTimeBasedOnProximity(utcDateStr);
		el.textContent = localDateFormatted;
	});
}

// Run the conversion when the document loads
window.addEventListener('DOMContentLoaded', convertUTCtoLocal);
// Also trigger the conversion after any HTMX request has settled (i.e., after content is swapped)
window.addEventListener('htmx:afterSettle', convertUTCtoLocal);

// Paste to Markdown on textarea
document.addEventListener("DOMContentLoaded", function() {
	const turndownService = new TurndownService({
		headingStyle: "atx",
		bulletListMarker: "-",
		codeBlockStyle: "fenced",
		emDelimiter: "*",
		strongDelimiter: "**",
		linkStyle: "inlined",
	});

	// Function to initialize markdown conversion for a textarea
	function initializeMarkdownConversion(textarea) {
		if (textarea.dataset.processed === "true") {
			return; // Prevent multiple initializations
		}
		textarea.dataset.processed = "true";

		// Handle text input to convert to HTML (optional, you may remove this if you only want paste handling)
		textarea.addEventListener("input", function() {
			const markdownText = textarea.value;
			const html = turndownService.turndown(markdownText);
			console.log(html); // You can replace this with whatever you need to do with the HTML, like rendering it or sending it to the server
		});

		// Handle paste event
		textarea.addEventListener("paste", function(event) {
			event.preventDefault(); // Prevent default paste action

			const clipboardData = (event.clipboardData || window.clipboardData);
			const pastedData = clipboardData.getData('text/html') || clipboardData.getData('text/plain');

			// Convert the pasted HTML or plain text to Markdown
			let markdown = pastedData;

			// If HTML is pasted, convert it
			if (clipboardData.getData('text/html')) {
				markdown = turndownService.turndown(pastedData);
			}

			// Insert the Markdown at the cursor position using execCommand for undo support
			insertMarkdownAtCursor(textarea, markdown);
		});
	}

	// Function to insert text at the cursor position while maintaining undo functionality
	function insertMarkdownAtCursor(textarea, markdown) {
		textarea.focus(); // Focus the textarea
		const selection = window.getSelection();
		if (document.queryCommandSupported("insertText")) {
			// Use the execCommand for inserting text to keep it in the undo stack
			document.execCommand("insertText", false, markdown);
		} else {
			// Fallback for browsers not supporting execCommand
			const start = textarea.selectionStart;
			const end = textarea.selectionEnd;

			const before = textarea.value.substring(0, start);
			const after = textarea.value.substring(end);

			// Set the new value with the markdown inserted
			textarea.value = before + markdown + after;

			// Set the new cursor position after the inserted markdown
			const newPosition = start + markdown.length;
			textarea.setSelectionRange(newPosition, newPosition);
		}
	}

	// Select and initialize any existing textareas
	document.querySelectorAll(".markdown-textarea").forEach(initializeMarkdownConversion);

	// MutationObserver to monitor DOM changes for new textareas
	const observer = new MutationObserver((mutations) => {
		mutations.forEach((mutation) => {
			if (mutation.type === "childList" && mutation.addedNodes.length > 0) {
				mutation.addedNodes.forEach((node) => {
					if (node.nodeType === Node.ELEMENT_NODE) {
						// Check for new textareas with the markdown-textarea class
						if (node.matches(".markdown-textarea")) {
							initializeMarkdownConversion(node);
						}

						// If new children are added that may contain textareas
						node.querySelectorAll(".markdown-textarea").forEach(initializeMarkdownConversion);
					}
				});
			}
		});
	});

	// Start observing the document body for changes
	observer.observe(document.body, { childList: true, subtree: true });
});
