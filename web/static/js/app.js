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
