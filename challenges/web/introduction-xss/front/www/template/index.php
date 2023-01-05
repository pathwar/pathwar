
<link href='/static/style.css' rel='stylesheet' type='text/css'>
<script src="static/script.js"></script>

<head>
    <meta charset="UTF-8">
    <title>Training XSS</title>
</head>
<body>
    <p>Les attaques XSS (Cross-Site Scripting) sont des injections qui contiennent du code malveillant.
    Ces failles peuvent apparaître dès qu'une application web utilise une entrée utilisateur.</p>

    <p>Notre entrée utilisateur peut alors modifier le comportement initialement prévu.
    Par exemple si l'on injecte du code de cette manière &lt;script&gt;code&lt;/script&gt; et que l'entrée n'est pas sécurisée on peut alors exécuter du code Javascript.</p>

    <p>
    Pour valider ce challenge, essayer d'exécuter le code alert ("XSS) à partir du formulaire ci-dessous.
    </p>
    <div class="form__group field">
        <input type="text" name="payload" id="payload" placeholder="Could you get the flag ? ...">
        <button onclick="sendRequest(document.getElementById('payload'))">GO</button>
    </div>
</body>