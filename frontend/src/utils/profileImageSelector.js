const common_path = '/icons/dashboard-icons/profile-icons';

const profileIcons = [
  `${common_path}/reshot-icon-cool-man-9GJ8R7456B.svg`,
  `${common_path}/reshot-icon-girl-9STZHYFACB.svg`,
  `${common_path}/reshot-icon-hairless-old-man-JE39NAM24U.svg`,
  `${common_path}/reshot-icon-happy-man-6CETQB5K7N.svg`,
  `${common_path}/reshot-icon-happy-woman-9G2BRLUEQF.svg`,
  `${common_path}/reshot-icon-happy-young-man-88H7Y6V293.svg`,
  `${common_path}/reshot-icon-happy-young-woman-VJG2UXYPDA.svg`,
  `${common_path}/reshot-icon-man-with-glasses-E3JYAV94TM.svg`,
  `${common_path}/reshot-icon-man-with-turtleneck-sweater-YMKJZTQ8R2.svg`,
  `${common_path}/reshot-icon-nerd-man-79HGAXCLRY.svg`,
  `${common_path}/reshot-icon-office-woman-BDJVUL8M57.svg`,
  `${common_path}/reshot-icon-old-bearded-man-DPM3VZTGVWB.svg`,
  `${common_path}/reshot-icon-relaxing-woman-JKBH2QDGNV.svg`,
  `${common_path}/reshot-icon-religious-man-MLPV28SHBE.svg`,
  `${common_path}/reshot-icon-serious-woman-5VREFYC54J.svg`,
  `${common_path}/reshot-icon-skinny-old-man-GWA8L32QVR.svg`,
  `${common_path}/reshot-icon-smiling-black-man-YV65U7GPMZ.svg`,
  `${common_path}/reshot-icon-smiling-man-E9RAGBZYSY.svg`,
  `${common_path}/reshot-icon-smiling-young-man-PGXTUZ6HML.svg`,
  `${common_path}/reshot-icon-woman-GLJT95XHRZ.svg`,
  `${common_path}/reshot-icon-worried-man-PBHTNEW9L.svg`,
  `${common_path}/reshot-icon-young-man-DK5RZF6X8JQ.svg`,
  `${common_path}/reshot-icon-young-person-ZMU2WGRB5K.svg`,
];
function getRandomProfileIcon() {
  const randomIndex = Math.floor(Math.random() * profileIcons.length);
  return profileIcons[randomIndex];
}
const randomIcon = getRandomProfileIcon();

export default randomIcon;
