# 📝 Customization Guide

This document lists all placeholders that need to be customized before using this project in production.

## 🔧 Required Customizations

### 1. LICENSE File

Replace these placeholders in `LICENSE`:

```
[your-email@example.com]        → Your contact email
[sales-email@example.com]       → Your sales/commercial email
[support-email@example.com]     → Your support email
[your-repo]                     → Your GitHub username/organization
```

**Lines to update:**
- Line 68: `Email: [your-email@example.com]`
- Line 69: `Web: https://github.com/[your-repo]/dockershield`
- Line 157: `- GPL-3.0 compliance: [your-email@example.com]`
- Line 158: `- Commercial licensing: [sales-email@example.com]`
- Line 159: `- General inquiries: [support-email@example.com]`
- Line 161: `Project Repository: https://github.com/[your-repo]/dockershield`

### 2. LICENSE-COMMERCIAL File

Replace these placeholders in `LICENSE-COMMERCIAL`:

```
[Your Address]                  → Your business address
[Your Email]                    → Your business email
[support-email@example.com]     → Your support email
[sales-email@example.com]       → Your sales email
[your-repo]                     → Your GitHub username/organization
[YOUR JURISDICTION]             → Legal jurisdiction (e.g., "France", "Germany")
[YOUR CITY]                     → City for arbitration
[Your Phone Number]             → Your contact phone
[Your Bank Information]         → Your payment details
```

**Lines to update:**
- Line 13-15: Licensor contact information
- Line 104: Support email
- Line 105: Issue tracker URL
- Line 203: Governing law jurisdiction
- Line 205: Arbitration city
- Line 261: Sales email
- Line 262: Repository URL
- Line 263: Phone number
- Line 256: Bank details

### 3. README.md

The README already contains the correct structure, but verify:

```
- Repository URLs point to your GitHub
- Contact emails are updated
- Pricing information matches your strategy
```

## 🚀 Quick Customization Script

You can use this bash script to replace all placeholders at once:

```bash
#!/bin/bash

# Set your information
YOUR_EMAIL="your.email@company.com"
SALES_EMAIL="sales@company.com"
SUPPORT_EMAIL="support@company.com"
YOUR_REPO="yourusername"
YOUR_ADDRESS="123 Your Street, Your City, Your Country"
YOUR_JURISDICTION="France"
YOUR_CITY="Paris"
YOUR_PHONE="+33 X XX XX XX XX"
YOUR_BANK="IBAN: FRXX XXXX XXXX XXXX XXXX XXXX XXX / BIC: XXXXXXXX"

# Replace in LICENSE
sed -i "s|\[your-email@example.com\]|$YOUR_EMAIL|g" LICENSE
sed -i "s|\[sales-email@example.com\]|$SALES_EMAIL|g" LICENSE
sed -i "s|\[support-email@example.com\]|$SUPPORT_EMAIL|g" LICENSE
sed -i "s|\[your-repo\]|$YOUR_REPO|g" LICENSE

# Replace in LICENSE-COMMERCIAL
sed -i "s|\[Your Address\]|$YOUR_ADDRESS|g" LICENSE-COMMERCIAL
sed -i "s|\[Your Email\]|$YOUR_EMAIL|g" LICENSE-COMMERCIAL
sed -i "s|\[support-email@example.com\]|$SUPPORT_EMAIL|g" LICENSE-COMMERCIAL
sed -i "s|\[sales-email@example.com\]|$SALES_EMAIL|g" LICENSE-COMMERCIAL
sed -i "s|\[your-repo\]|$YOUR_REPO|g" LICENSE-COMMERCIAL
sed -i "s|\[YOUR JURISDICTION\]|$YOUR_JURISDICTION|g" LICENSE-COMMERCIAL
sed -i "s|\[YOUR CITY\]|$YOUR_CITY|g" LICENSE-COMMERCIAL
sed -i "s|\[Your Phone Number\]|$YOUR_PHONE|g" LICENSE-COMMERCIAL
sed -i "s|\[Your Bank Information\]|$YOUR_BANK|g" LICENSE-COMMERCIAL

echo "✅ All placeholders replaced successfully!"
```

Save this as `customize.sh`, update the variables at the top, make it executable with `chmod +x customize.sh`, and run it with `./customize.sh`.

## ✅ Verification Checklist

After customization, verify:

- [ ] All `[placeholder]` values are replaced
- [ ] Email addresses are valid and monitored
- [ ] Repository URL is correct
- [ ] Pricing tiers match your business model
- [ ] Legal jurisdiction is appropriate
- [ ] Bank/payment information is accurate
- [ ] Phone numbers are reachable

## 📋 Optional Customizations

### Pricing Strategy

You may want to adjust the pricing in `LICENSE-COMMERCIAL`:

- **Startup**: Currently €500/year (< 10 employees)
- **SME**: Currently €2,000/year (< 100 employees)
- **Enterprise**: Currently €10,000/year (> 100 employees)
- **OEM**: Custom pricing

### Support Levels

Adjust response times based on your capacity:

- Startup: 48h response time
- SME: 24h response time
- Enterprise: 4h response time

### Payment Terms

Consider your preferred:
- Payment schedule (annual/quarterly/one-time)
- Late payment interest (currently 1.5%/month)
- Grace periods
- Refund policy

## 🔒 Security Considerations

Before going to production:

1. **Review default security filters** in `config/defaults.go`
2. **Set appropriate PROXY_CONTAINER_NAME** environment variable
3. **Configure network isolation** with PROXY_NETWORK_NAME
4. **Enable advanced filters** if needed
5. **Audit access rules** for your use case

## 📞 Support

For questions about customization:
- Documentation: See README.md, docs/SECURITY.md, docs/ENV_FILTERS.md
- Issues: Report at [your repository]/issues
- Commercial inquiries: [your sales email]
