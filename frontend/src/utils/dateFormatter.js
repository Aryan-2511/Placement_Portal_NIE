export default function dateFormatter(dateString, type) {
  let formattedDate = new Date(dateString);
  if (type === 'string') {
    formattedDate = formattedDate.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  } else if (type === 'date') {
    formattedDate = formattedDate.toISOString().split('T')[0];
  }

  return formattedDate;
}
