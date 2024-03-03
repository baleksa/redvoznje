function toCSV(inputArray, separator = ",") {
	var header = "sat,radni dan,subota,nedelja\n"
	let rowsAsString = inputArray.map(row => {
		return row.join(separator);
	})
	var csvFormat = rowsAsString.join("\n");

	csvFormat = header + csvFormat
	return csvFormat
}

function makeTextFile(text) {
	var data = new Blob([text], {
		type:
			'text/plain'
	});

	var textFile = window.URL.createObjectURL(data);

	return textFile;
};

function downloadAllSchedules() {
	$.fn.dataTable.tables().forEach(table => {
		var dtable = $(table).DataTable()
		var data = dtable.data()

		var link = document.createElement("a");
		link.download = `${getTitle(dtable.table())}.csv`;
		link.href = makeTextFile(toCSV(data));
		link.click();
	});

}

function evalValues() {
	var values = []

	$.fn.DataTable.tables().forEach(table => {
		var dtable = $(table).DataTable()
		var value = dtable.data()
		values.push(value)
	})
	return values
}

function getTitle(table) {
	var rawTitle = table.container().parentElement.parentElement.parentElement.querySelector("p").textContent
	return rawTitle.replace("\n", " ")
}

downloadAllSchedules()
