// Theme exports
export {
  theme,
  darkTheme,
  spacing,
  borderRadius,
  customColors,
} from './theme';
export type { AppTheme, CustomColors } from './theme';

// Component exports
export * from './components';

// Re-export commonly used React Native Paper components
export {
  Appbar,
  Avatar,
  Badge,
  Banner,
  Chip,
  DataTable,
  Divider,
  FAB,
  IconButton,
  List,
  Menu,
  ProgressBar,
  RadioButton,
  Searchbar,
  SegmentedButtons,
  Snackbar,
  Surface,
  Switch,
  Text,
  Tooltip,
  TouchableRipple,
  useTheme,
  Provider as PaperProvider,
} from 'react-native-paper';
