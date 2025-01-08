export default function dateFormatter(dateString) {
  let formattedDate = new Date(dateString);
  formattedDate = formattedDate.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  return formattedDate;
}
