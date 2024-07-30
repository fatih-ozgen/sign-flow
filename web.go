package main

import (
	"html/template"
	"net/http"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>User Management</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 500px;
            margin: 0 auto;
            padding: 20px;
        }
        form {
            display: flex;
            flex-direction: column;
        }
        input, button {
            margin: 10px 0;
            padding: 10px;
        }
        button {
            cursor: pointer;
        }
        #googleSignUp {
            background-color: #4285F4;
            color: white;
            border: none;
        }
        #submitButton {
            background-color: #4CAF50;
            color: white;
            border: none;
        }
        #cancelButton {
            background-color: #f44336;
            color: white;
            border: none;
        }
    </style>
</head>
<body>
    <h1>User Management</h1>
    <form id="signupForm">
        <button type="button" id="googleSignUp">Sign up with Google</button>
        <input type="text" id="username" placeholder="Username" required>
        <div style="position: relative;">
            <input type="password" id="password" placeholder="Password" required>
            <button type="button" id="showPassword" style="position: absolute; right: 5px; top: 50%; transform: translateY(-50%);">Show</button>
        </div>
        <button type="submit" id="submitButton">Create User</button>
        <button type="button" id="cancelButton">Cancel</button>
    </form>

    <script>
        document.getElementById('googleSignUp').addEventListener('click', function() {
            window.location.href = '/auth/google/login';
        });

        document.getElementById('showPassword').addEventListener('click', function() {
            var passwordInput = document.getElementById('password');
            if (passwordInput.type === 'password') {
                passwordInput.type = 'text';
                this.textContent = 'Hide';
            } else {
                passwordInput.type = 'password';
                this.textContent = 'Show';
            }
        });

        document.getElementById('cancelButton').addEventListener('click', function() {
            document.getElementById('signupForm').reset();
        });

        document.getElementById('signupForm').addEventListener('submit', function(e) {
            e.preventDefault();
            var username = document.getElementById('username').value;
            var password = document.getElementById('password').value;
            
            fetch('/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({username: username, password: password}),
            })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
                if (data.membership_id) {
                    alert('Your membership ID is: ' + data.membership_id);
                }
            })
            .catch((error) => {
                console.error('Error:', error);
                alert('An error occurred while creating the user.');
            });
        });
    </script>
</body>
</html>
`

func serveWebPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("webpage").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
