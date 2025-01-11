export const verifyLogin = (email, password) => {
  const emailRegex = /^(\d{4})(cs|is|ec|ee|me|cv)_[a-z]+_(a|b)@nie\.ac\.in$/i;
  const passwordRegex = /^(?=.*[A-Z])(?=.*\d)[A-Za-z\d]{8,}$/;

  // Validate email
  if (!emailRegex.test(email)) {
    return {
      valid: false,
      error:
        "Invalid email format. Please follow the specified format: '20XXbb_johndoe_s@nie.ac.in'.",
    };
  }

  // Validate password
  // if (!passwordRegex.test(password)) {
  //   return {
  //     valid: false,
  //     error:
  //       'Password must be at least 8 characters long, with at least one uppercase letter and one number.',
  //   };
  // }

  // If both are valid
  return { valid: true };
};
