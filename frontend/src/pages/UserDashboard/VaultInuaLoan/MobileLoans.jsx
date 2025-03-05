import React, { useState } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Grid, 
  Card, 
  CardContent, 
  Button, 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  TextField,
  Paper,
  Divider
} from '@mui/material';
import { PhoneAndroid, MonetizationOn, Timeline } from '@mui/icons-material';

const MobileLoans = () => {
  const [openLoanDialog, setOpenLoanDialog] = useState(false);
  const [loanAmount, setLoanAmount] = useState('');
  const [loanPurpose, setLoanPurpose] = useState('');

  const loanTypes = [
    {
      title: 'Instant Mobile Loan',
      description: 'Quick loan directly to your mobile wallet',
      icon: <PhoneAndroid />,
      maxLimit: 30000,
      interestRate: 15
    },
    {
      title: 'Flex Mobile Loan',
      description: 'Flexible mobile financing with easy repayment',
      icon: <MonetizationOn />,
      maxLimit: 75000,
      interestRate: 12
    },
    {
      title: 'Extended Mobile Loan',
      description: 'Longer-term mobile financial support',
      icon: <Timeline />,
      maxLimit: 150000,
      interestRate: 10
    }
  ];

  const handleOpenLoanDialog = () => {
    setOpenLoanDialog(true);
  };

  const handleCloseLoanDialog = () => {
    setOpenLoanDialog(false);
  };

  const handleApplyLoan = () => {
    // TODO: Implement loan application logic
    console.log('Mobile Loan Application', { loanAmount, loanPurpose });
    handleCloseLoanDialog();
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Mobile Loans
      </Typography>
      
      <Grid container spacing={4}>
        {loanTypes.map((loan, index) => (
          <Grid item xs={12} md={4} key={index}>
            <Card 
              sx={{ 
                height: '100%', 
                display: 'flex', 
                flexDirection: 'column',
                transition: 'transform 0.3s',
                '&:hover': {
                  transform: 'scale(1.05)'
                }
              }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  {loan.icon}
                  <Typography variant="h6" sx={{ ml: 2 }}>
                    {loan.title}
                  </Typography>
                </Box>
                <Typography variant="body2" color="text.secondary">
                  {loan.description}
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                  <Typography variant="body2">
                    Max Limit: Ksh {loan.maxLimit.toLocaleString()}
                  </Typography>
                  <Typography variant="body2">
                    Interest: {loan.interestRate}%
                  </Typography>
                </Box>
              </CardContent>
              <Button 
                variant="contained" 
                color="primary" 
                fullWidth
                onClick={handleOpenLoanDialog}
              >
                Apply Now
              </Button>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Dialog open={openLoanDialog} onClose={handleCloseLoanDialog}>
        <DialogTitle>Apply for Mobile Loan</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Loan Amount"
            type="number"
            fullWidth
            value={loanAmount}
            onChange={(e) => setLoanAmount(e.target.value)}
          />
          <TextField
            margin="dense"
            label="Loan Purpose"
            fullWidth
            multiline
            rows={3}
            value={loanPurpose}
            onChange={(e) => setLoanPurpose(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseLoanDialog} color="primary">
            Cancel
          </Button>
          <Button onClick={handleApplyLoan} color="primary" variant="contained">
            Submit Application
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default MobileLoans;
